package downloader

import (
	"os"
	s3client "service-worker-sqs-s3-postgres/dataproviders/awss3"
	"service-worker-sqs-s3-postgres/dataproviders/utils"

	"github.com/labstack/gommon/log"
	"github.com/spf13/afero"
	"go.uber.org/zap"
)

// S3Downloader represents a S3 handler.
type S3Downloader struct {
	fs  afero.Fs
	s3  *s3client.ClientS3
	log *zap.SugaredLogger
}

// NewDownloader instances a S3Downloader.
func NewDownloader(s3 *s3client.ClientS3, log *zap.SugaredLogger) (*S3Downloader, error) {
	fs := afero.NewOsFs()
	if err := fs.MkdirAll(utils.TmpPath, os.ModePerm); err != nil {
		return nil, err
	}
	return &S3Downloader{fs: fs, s3: s3, log: log}, nil
}

// Download downloads a file from S3 bucket.
func (d *S3Downloader) Download(bucket, key string) (string, error) {
	localPath := utils.CreateLocalFileName(key)
	dst, err := d.fs.Create(localPath)
	if err != nil {
		log.Errorf("s3downloader: error creating dst file %s: %v", localPath, err)
		return "", err
	}
	defer utils.Close(dst, d.log)

	err = d.s3.DownloadFile(bucket, key, dst)
	if err != nil {
		d.log.Errorf("s3downloader: error downloading file %s. %v", key, err)
		return "", err
	}
	return localPath, nil
}

// Delete local file.
func (d *S3Downloader) Delete(file string) error {
	return d.fs.Remove(file)
}
