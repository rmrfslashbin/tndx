package storage

import (
	"bytes"
	"compress/gzip"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type S3Option func(c *S3StorageDriver)

type S3StorageDriver struct {
	driverName string
	s3Bucket   string
	s3Region   string
}

func NewS3Storage(opts ...func(*S3StorageDriver)) *S3StorageDriver {
	config := &S3StorageDriver{}
	config.driverName = "S3"

	// apply the list of options to Config
	for _, opt := range opts {
		opt(config)
	}
	return config
}

func SetS3Bucket(s3Bucket string) S3Option {
	return func(config *S3StorageDriver) {
		config.s3Bucket = s3Bucket
	}
}

func SetS3Region(s3Region string) S3Option {
	return func(config *S3StorageDriver) {
		config.s3Region = s3Region
	}
}

// PutObject uploads data to an S3 bucket.
func (config *S3StorageDriver) Put(key string, body []byte) error {
	// *s3manager.UploadOutput

	// The session the S3 Uploader will use
	// Specify profile for config and region for requests.
	s3Session := session.Must(session.NewSessionWithOptions(session.Options{
		Config: aws.Config{Region: aws.String(config.s3Region)},
	}))

	// Create an uploader with the session and default options.
	uploader := s3manager.NewUploader(s3Session)

	// gzip data
	var buf bytes.Buffer
	zw := gzip.NewWriter(&buf)
	_, err := zw.Write(body)
	if err != nil {
		return err
	}

	// Bail out if we got an error while compressing.
	if err := zw.Close(); err != nil {
		return err
	}

	// Append ".gz" to the key (filename).
	key = key + ".gz"

	// Upload input parameters
	upParams := &s3manager.UploadInput{
		Bucket: &config.s3Bucket,
		Key:    &key,
		Body:   &buf,
	}

	// Perform an upload.
	_, err = uploader.Upload(upParams)
	return err
	// return result, err
}

func (config *S3StorageDriver) GetDriverName() string {
	return config.driverName
}
