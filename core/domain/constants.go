package domain

type Session string

const (
	SQS Session = "sqs"
	S3  Session = "s3"
)
