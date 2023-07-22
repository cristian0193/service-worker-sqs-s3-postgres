package utils

import (
	"fmt"
	"go.uber.org/zap"
	"io"
	"path"
	"strconv"
	"time"
)

const (
	TmpPath = "/tmp/files/s3"
	Layout  = "2006-01-02T15-04-05.000000000"
)

// CreateLocalFileName return path to local.
func CreateLocalFileName(key string) string {
	fileName := AddTimestamp(key, time.Now())
	localPath := path.Join(TmpPath, fileName)
	return localPath
}

// AddTimestamp returns the read time in filename.
func AddTimestamp(filename string, time time.Time) string {
	originalFileName := path.Base(filename)
	ft := time.Format(Layout)
	return fmt.Sprintf("%s_%s", originalFileName, ft)
}

// Close closes the file.
func Close(c io.Closer, log *zap.SugaredLogger) {
	if err := c.Close(); err != nil {
		log.Error(err)
	}
}

// StringToInt64 convert string to int64.
func StringToInt64(s string) *int64 {
	if s == "" {
		return nil
	}
	n, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return nil
	}
	return &n
}