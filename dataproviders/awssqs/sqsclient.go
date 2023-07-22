package awssqs

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/aws/aws-sdk-go/service/sqs/sqsiface"
)

// ClientSQS represents SQS client.
type ClientSQS struct {
	api               sqsiface.SQSAPI
	url               string
	maxMessages       int64
	visibilityTimeout int64
}

// NewSQSClient instances of a Client to connect SQS with session as parameter.
func NewSQSClient(sess *session.Session, url string, maxMessages, visibilityTimeout int) (*ClientSQS, error) {
	return &ClientSQS{
		api:               sqs.New(sess),
		url:               url,
		maxMessages:       int64(maxMessages),
		visibilityTimeout: int64(visibilityTimeout),
	}, nil
}

// GetMessages retrieves messages from SQS.
func (s *ClientSQS) GetMessages() ([]*sqs.Message, error) {
	params := &sqs.ReceiveMessageInput{
		QueueUrl:            aws.String(s.url),
		MaxNumberOfMessages: aws.Int64(s.maxMessages),
		AttributeNames: []*string{
			aws.String("All"),
		},
		MessageAttributeNames: []*string{
			aws.String("All"),
		},
		WaitTimeSeconds:   aws.Int64(20),
		VisibilityTimeout: aws.Int64(s.visibilityTimeout),
	}

	res, err := s.api.ReceiveMessage(params)
	if err != nil {
		return nil, err
	}

	return res.Messages, nil
}

// DeleteMessage deletes messages from SQS.
func (s *ClientSQS) DeleteMessage(msg *sqs.Message) error {
	params := &sqs.DeleteMessageInput{
		QueueUrl:      aws.String(s.url),
		ReceiptHandle: msg.ReceiptHandle,
	}
	_, err := s.api.DeleteMessage(params)

	return err
}
