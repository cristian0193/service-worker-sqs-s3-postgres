package builder

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws/session"
	"go.uber.org/zap"
	"service-worker-sqs-s3-postgres/core/domain"
	"service-worker-sqs-s3-postgres/dataproviders/awss3"
	"service-worker-sqs-s3-postgres/dataproviders/awss3/downloader"
	"service-worker-sqs-s3-postgres/dataproviders/awssqs"
	"service-worker-sqs-s3-postgres/dataproviders/consumer"
	rfiledata "service-worker-sqs-s3-postgres/dataproviders/postgres/repository/filedata"
	rmetadata "service-worker-sqs-s3-postgres/dataproviders/postgres/repository/metadata"
)

// NewConsumer define all usecases to instantiate SQS.
func NewConsumer(logger *zap.SugaredLogger,
	config *Configuration,
	sessSQS *session.Session,
	sessS3 *session.Session,
	rfd rfiledata.IFileDataRepository,
	rmd rmetadata.IMetaDataRepository) (domain.Source, error) {

	s3, err := awss3.NewS3Client(sessS3, config.S3Bucket)
	if err != nil {
		return nil, fmt.Errorf("error awssqs.NewSQSClient: %w", err)
	}

	sqs, err := awssqs.NewSQSClient(sessSQS, config.SQSUrl, config.SQSMaxMessages, config.SQSVisibilityTimeout)
	if err != nil {
		return nil, fmt.Errorf("error awssqs.NewSQSClient: %w", err)
	}

	download, err := downloader.NewDownloader(s3, logger)
	if err != nil {
		return nil, fmt.Errorf("error downloader.NewDownloader: %w", err)
	}

	source, err := consumer.New(sqs, download, logger, config.SQSMaxMessages, rfd, rmd)
	if err != nil {
		return nil, fmt.Errorf("error consumer.New: %w", err)
	}

	return source, nil
}
