package builder

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"service-worker-sqs-s3-postgres/core/domain"
)

// NewSession define all configuration to instantiate a session aws.
func NewSession(config *Configuration, typeSession domain.Session) (*session.Session, error) {
	sessionConfig := &aws.Config{
		Region:      aws.String(config.Region),
		Credentials: credentials.NewStaticCredentials(config.AccessKey, config.SecretKey, ""),
	}

	switch typeSession {
	case domain.SQS:
		sessionConfig.Endpoint = aws.String(config.SQSUrl)
		sessionConfig.MaxRetries = aws.Int(3)
		break
	case domain.S3:
		sessionConfig.MaxRetries = aws.Int(5)
		break
	}

	sess, err := session.NewSession(sessionConfig)
	if err != nil {
		return nil, err
	}

	return session.Must(sess, err), nil
}
