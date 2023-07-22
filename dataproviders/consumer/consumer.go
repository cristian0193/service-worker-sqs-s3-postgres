package consumer

import (
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/tidwall/gjson"
	"go.uber.org/zap"
	"service-worker-sqs-s3-postgres/core/domain"
	"service-worker-sqs-s3-postgres/dataproviders/awss3/downloader"
	"service-worker-sqs-s3-postgres/dataproviders/awssqs"
	repository "service-worker-sqs-s3-postgres/dataproviders/postgres/repository/events"
	"sync"
)

// SQSSource event stream representation to SQS.
type SQSSource struct {
	sqs         *awssqs.ClientSQS
	download    *downloader.S3Downloader
	log         *zap.SugaredLogger
	maxMessages int
	closed      bool
	repo        repository.IEventRepository
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
func New(sqsClient *awssqs.ClientSQS, download *downloader.S3Downloader, logger *zap.SugaredLogger, maxMessages int, repo repository.IEventRepository) (*SQSSource, error) {
	return &SQSSource{
		sqs:         sqsClient,
		download:    download,
		log:         logger,
		maxMessages: maxMessages,
		repo:        repo,
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

	logger.Infof("Step 2.1 - Event from path: %s", filename)

	//TODO: read file and mapping in struct

	//TODO: saved in database

	/*eventDB := &entity.Events{
		ID:      *msg.MessageId,
		Message: records.Message,
		Date:    time.Now().Format(time.RFC3339),
	}

	if err = s.repo.Insert(eventDB); err != nil {
		logger.Errorf("Error inserting message: %v", err)
	}

	logger.Info("Step 3 - File saved in postgres")*/

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
	logger.Infof("Step 4 - Event produced for ID = %s)", trackID)
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
		logger.Infof("Successful deleted sqs message")
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
