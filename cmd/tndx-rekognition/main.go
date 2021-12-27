package main

import (
	"context"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/service/rekognition/types"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/rmrfslashbin/tndx/pkg/rekognition"
	"github.com/sirupsen/logrus"
)

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

func handler(ctx context.Context, event events.S3Event) error {
	log.WithFields(logrus.Fields{
		"event": event,
	}).Info("received message")

	for _, record := range event.Records {

		rk := rekognition.NewImageProcessor(
			rekognition.SetRegion(aws_region),
			rekognition.SetLogger(log),
		)

		output, err := rk.Process(&types.S3Object{
			Bucket: aws.String(record.S3.Bucket.Name),
			Name:   aws.String(record.S3.Object.Key),
			//Version: aws.String(record.S3.Object.VersionID),
		})
		if err != nil {
			log.WithFields(logrus.Fields{
				"error": err,
			}).Error("error processing image")
			return err
		}
		log.WithFields(logrus.Fields{
			"output": output,
		}).Info("image processed")
	}
	return nil
}
