package builder

import (
	"go.uber.org/zap"
	"service-worker-sqs-s3-postgres/core/domain"
	"service-worker-sqs-s3-postgres/dataproviders/processor"
)

// NewProcessor define all usecases to be instantiated Processor associated with the consumer.
func NewProcessor(logger *zap.SugaredLogger, source domain.Source) (*processor.Processor, error) {
	return processor.New(logger, source)
}
