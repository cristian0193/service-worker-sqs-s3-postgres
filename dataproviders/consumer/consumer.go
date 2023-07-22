package consumer

import (
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/tidwall/gjson"
	"go.uber.org/zap"
	"service-worker-sqs-s3-postgres/core/domain"
	"service-worker-sqs-s3-postgres/core/domain/entity"
	"service-worker-sqs-s3-postgres/dataproviders/awss3/downloader"
	"service-worker-sqs-s3-postgres/dataproviders/awssqs"
	"service-worker-sqs-s3-postgres/dataproviders/consumer/csvreader"
	"service-worker-sqs-s3-postgres/dataproviders/mapper"
	rfiledata "service-worker-sqs-s3-postgres/dataproviders/postgres/repository/filedata"
	rmetadata "service-worker-sqs-s3-postgres/dataproviders/postgres/repository/metadata"
	"sync"
)

// SQSSource event stream representation to SQS.
type SQSSource struct {
	sqs         *awssqs.ClientSQS
	download    *downloader.S3Downloader
	log         *zap.SugaredLogger
	maxMessages int
	closed      bool
	rFiledata   rfiledata.IFileDataRepository
	rMetadata   rmetadata.IMetaDataRepository
	wg          sync.WaitGroup
}

type s3Event struct {
	bucket     string
	key        string
	fileSize   int64
	sqsMessage *sqs.Message
}

var (
	errInvalidJSON        = errors.New("invalid json")
	errNoRecordsFound     = errors.New("no records found")
	errInvalidEventSource = errors.New("invalid event source")
)

// New return an event stream instance from SQS.
func New(sqsClient *awssqs.ClientSQS, download *downloader.S3Downloader, logger *zap.SugaredLogger, maxMessages int, rfd rfiledata.IFileDataRepository, rmd rmetadata.IMetaDataRepository) (*SQSSource, error) {
	return &SQSSource{
		sqs:         sqsClient,
		download:    download,
		log:         logger,
		maxMessages: maxMessages,
		rFiledata:   rfd,
		rMetadata:   rmd,
		wg:          sync.WaitGroup{},
	}, nil
}

// Consume opens a channel and sends entities created from SQS messages.
func (s *SQSSource) Consume() <-chan *domain.Event {
	out := make(chan *domain.Event, s.maxMessages)
	go func() {
		for {
			if s.closed {
				break
			}
			messages, err := s.sqs.GetMessages()
			if err != nil {
				s.log.Errorf("Error getting messages from SQS: %v", err)
				continue
			}
			if len(messages) == 0 {
				s.log.Debug("No messages found from SQS")
			}
			for _, msg := range messages {
				s.processMessage(msg, out)
			}
			s.wg.Wait()
		}
		close(out)
	}()

	return out
}

// processMessage read message in queue.
func (s *SQSSource) processMessage(msg *sqs.Message, out chan *domain.Event) {
	trackID := createTrackID(msg)
	logger := s.log.With("trackId", trackID)

	logger.Infof("Step 1 - Start to process SQS event")

	s3Event, err := toS3Event(msg)
	if err != nil {
		logger.Errorf("Error processing message from SQS: %v", err)
		if err = s.sqs.DeleteMessage(msg); err != nil {
			logger.Errorf("Error deleting message from SQS: %v", err)
		}
		return
	}

	logger.Info("Step 2 - Starts the process of downloading the file from S3")

	filename, err := s.download.Download(s3Event.bucket, s3Event.key)
	if err != nil {
		logger.Errorf("Error processing message from SQS in [path = %s]: %v", s3Event.key, err)
		if err = s.sqs.DeleteMessage(msg); err != nil {
			logger.Errorf("Error deleting message from SQS in [path = %s]: %v", s3Event.key, err)
		}
		return
	}

	logger.Infof("Step 3 - Event from path: %s", filename)

	filedata, err := fileMapping(filename, logger)
	if err != nil {
		logger.Errorf("Error processing file from CSV in [path = %s]: %v", s3Event.key, err)
		if err = s.sqs.DeleteMessage(msg); err != nil {
			logger.Errorf("Error deleting message from SQS in [path = %s]: %v", s3Event.key, err)
		}
		return
	}

	if err = s.rFiledata.Insert(filedata); err != nil {
		logger.Errorf("Error inserting message in FileData: %v", err)
	}

	logger.Info("Step 4 - File saved in postgres: FileData")

	metadata := entity.MetaData{
		TrackID:  trackID,
		Bucket:   s3Event.bucket,
		FileName: filename,
		Key:      s3Event.key,
		Size:     s3Event.fileSize,
	}

	if err = s.rMetadata.Insert(metadata); err != nil {
		logger.Errorf("Error inserting message in MetaData: %v", err)
	}

	logger.Info("Step 4 - File saved in postgres: MetaData")

	event := &domain.Event{
		TrackID:       trackID,
		File:          s3Event.key,
		Bucket:        s3Event.bucket,
		OriginalEvent: s3Event,
		FileSize:      s3Event.fileSize,
		Log:           logger,
		Filename:      filename,
	}
	s.wg.Add(1)
	logger.Infof("Step 5 - Event produced for ID = %s)", trackID)
	out <- event
}

// Processed notify that event of consolidate file was processed.
func (s *SQSSource) Processed(event *domain.Event) error {
	defer s.wg.Done()
	logger := event.Log

	if err := s.download.Delete(event.Filename); err != nil {
		logger.Errorf("Error deleting local file: %v", err)
	}

	if s3Event, ok := event.OriginalEvent.(*s3Event); ok {
		if err := s.sqs.DeleteMessage(s3Event.sqsMessage); err != nil {
			logger.Errorf("Deleting of sqs message. %v", err)
			return err
		}
		logger.Infof("Step 6 - Successful deleted sqs message")
		return nil
	}
	logger.Warnf("Event isn't sqs message")
	return nil
}

// Close the event stream.
func (s *SQSSource) Close() error {
	s.closed = true
	s.wg.Wait()
	return nil
}

// ---------- Helpers ------------ //

func createTrackID(msg *sqs.Message) string {
	receiveCount := "0"
	val, ok := msg.Attributes[sqs.MessageSystemAttributeNameApproximateReceiveCount]
	if ok {
		receiveCount = *val
	}
	return fmt.Sprintf("%s-%s", *msg.MessageId, receiveCount)
}

func toS3Event(msg *sqs.Message) (*s3Event, error) {
	body := *msg.Body
	if !gjson.Valid(body) {
		return nil, errInvalidJSON
	}

	record := gjson.Get(body, "Records.0")
	if !record.Exists() {
		return nil, errNoRecordsFound
	}

	src := record.Get("eventSource")
	if !src.Exists() || src.String() != "aws:s3" {
		return nil, fmt.Errorf(`"%v": %w`, src, errInvalidEventSource)
	}

	return &s3Event{
		bucket:     record.Get("s3.bucket.name").String(),
		key:        record.Get("s3.object.key").String(),
		fileSize:   record.Get("s3.object.size").Int(),
		sqsMessage: msg,
	}, nil
}

func fileMapping(fileName string, logger *zap.SugaredLogger) ([]entity.FileData, error) {
	csv, err := csvreader.Read(fileName, logger)
	if err != nil {
		return nil, err
	}

	filedata := make([]entity.FileData, 0)
	for _, row := range csv {
		f := mapper.ToEntityFileData(&row)
		filedata = append(filedata, *f)
	}
	return filedata, nil
}
