package main

import (
	"bytes"
	"context"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

type Config struct {
	AwsRegion        string `yaml:"aws_region"`
	AwsProfile       string `yaml:"aws_profile"`
	DDBTablePrefix   string `yaml:"ddb_table_prefix"`
	TwitterAPIKey    string `yaml:"twitter_api_key"`
	TwitterAPISecret string `yaml:"twitter_api_secret"`
}

type S3Request struct {
	Bucket string `json:"bucket"`
	Key    string `json:"key"`
	Region string `json:"region"`
}

var (
	aws_region string
	log        *logrus.Logger
)

func init() {
	log = logrus.New()
	log.SetLevel(logrus.InfoLevel)
	log.SetFormatter(&logrus.JSONFormatter{})
	aws_region = os.Getenv("AWS_REGION")
}

func main() {
	// Catch errors
	var err error
	defer func() {
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"error": err,
			}).Fatal("main crashed")
		}
	}()
	lambda.Start(handler)
}

func handler(ctx context.Context, s3Request S3Request) error {
	cfg := Config{
		AwsRegion:        os.Getenv("REGION"),
		AwsProfile:       os.Getenv("PROFILE"),
		DDBTablePrefix:   os.Getenv("DDB_TABLE_PREFIX"),
		TwitterAPIKey:    os.Getenv("TWITTER_API_KEY"),
		TwitterAPISecret: os.Getenv("TWITTER_API_SECRET"),
	}
	yml, err := yaml.Marshal(&cfg)
	if err != nil {
		log.WithFields(
			logrus.Fields{
				"error":  err,
				"config": cfg,
				"s3":     s3Request,
			}).Error("failed to marshal config")
		return err
	}

	awsconfig, err := config.LoadDefaultConfig(context.TODO(), func(o *config.LoadOptions) error {
		o.Region = s3Request.Region
		return nil
	})
	if err != nil {
		log.WithFields(
			logrus.Fields{
				"error":  err,
				"config": cfg,
				"s3":     s3Request,
			}).Error("failed to config S3")
		return err
	}

	s3Svc := s3.NewFromConfig(awsconfig)
	obj, err := s3Svc.PutObject(ctx, &s3.PutObjectInput{
		Bucket: &s3Request.Bucket,
		Key:    &s3Request.Key,
		Body:   bytes.NewReader(yml),
	})
	if err != nil {
		log.WithFields(
			logrus.Fields{
				"error":  err,
				"config": cfg,
				"s3":     s3Request,
			}).Error("failed to put object")
		return err
	}
	log.WithFields(logrus.Fields{
		"config":       cfg,
		"s3":           s3Request,
		"outputObject": obj,
	}).Info("successfully wrote config")
	return nil
}
