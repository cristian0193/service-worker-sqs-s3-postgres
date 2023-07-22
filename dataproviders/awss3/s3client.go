package awss3

import (
	"io"
	"service-worker-sqs-s3-postgres/dataproviders/utils"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

// ClientS3 represents an AWS s3 client.
type ClientS3 struct {
	api        s3iface.S3API
	downloader *s3manager.Downloader
	bucket     string
}

// NewS3Client instantiates a new Client.
func NewS3Client(sess *session.Session, bucket string, cfgs ...*aws.Config) (*ClientS3, error) {
	api := s3.New(sess, cfgs...)
	return &ClientS3{
		api:        api,
		downloader: s3manager.NewDownloaderWithClient(api),
		bucket:     bucket,
	}, nil
}

// DownloadFile download a file from S3 bucket.
func (c *ClientS3) DownloadFile(bucket, key string, file io.WriterAt) error {
	if len(bucket) == 0 {
		bucket = c.bucket
	}

	params := &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}

	err := utils.Do(5, 3*time.Second, func() (bool, error) {
		_, err := c.downloader.Download(file, params)
		return true, err
	})
	if err != nil {
		return err
	}
	return nil
}
