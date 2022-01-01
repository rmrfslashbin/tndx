package storage

import (
	"bytes"
	"compress/gzip"
	"context"
	"io"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Option func(c *S3Storage)

type S3Storage struct {
	driverName string
	bucket     string
	region     string
	profile    string
	svc        *s3.Client
}

func NewS3Storage(opts ...func(*S3Storage)) *S3Storage {
	cfg := &S3Storage{}
	cfg.driverName = "S3"

	// apply the list of options to Config
	for _, opt := range opts {
		opt(cfg)
	}

	if cfg.region == "" {
		log.Fatal("region is required")
	}

	c, err := config.LoadDefaultConfig(context.TODO(), func(o *config.LoadOptions) error {
		o.Region = cfg.region
		if cfg.profile != "" {
			o.SharedConfigProfile = cfg.profile
		}
		return nil
	})
	if err != nil {
		panic(err)
	}

	cfg.svc = s3.NewFromConfig(c)

	return cfg
}

func SetProfile(profile string) S3Option {
	return func(config *S3Storage) {
		config.profile = profile
	}
}

func SetS3Bucket(bucket string) S3Option {
	return func(config *S3Storage) {
		config.bucket = bucket
	}
}

func SetS3Region(region string) S3Option {
	return func(config *S3Storage) {
		config.region = region
	}
}

func (config *S3Storage) Put(key string, body []byte) error {
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

	if _, err := config.svc.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(config.bucket),
		Key:    aws.String(key),
		Body:   bytes.NewReader(buf.Bytes()),
	}); err != nil {
		return err
	}
	return nil
}

func (config *S3Storage) PutStream(key string, fp io.Reader) error {
	if _, err := config.svc.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(config.bucket),
		Key:    aws.String(key),
		Body:   fp,
	}); err != nil {
		return err
	}
	return nil
}
