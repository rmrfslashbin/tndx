package main

import (
	"context"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/service/rekognition/types"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/rmrfslashbin/tndx/pkg/database"
	"github.com/rmrfslashbin/tndx/pkg/rekognition"
	"github.com/rmrfslashbin/tndx/pkg/ssmparams"
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

	params := ssmparams.NewSSMParams(
		ssmparams.SetRegion(aws_region),
		ssmparams.SetLogger(log),
	)

	outputs, err := params.GetParams([]string{
		os.Getenv("DDB_TABLE_PREFIX"),
	})
	if err != nil {
		log.WithFields(logrus.Fields{
			"error": err,
		}).Error("failed to get parameters for DDB_TABLE_PREFIX")
		return err
	}

	if len(outputs.InvalidParameters) > 0 {
		log.WithFields(logrus.Fields{
			"invalid_parameters": outputs.InvalidParameters,
		}).Error("invalid parameters for DDB_TABLE_PREFIX")
		return err
	}

	rk := rekognition.NewImageProcessor(
		rekognition.SetRegion(aws_region),
		rekognition.SetLogger(log),
	)

	ddb := database.NewDDB(
		database.SetDDBLogger(log),
		database.SetDDBTablePrefix(outputs.Params[os.Getenv("DDB_TABLE_PREFIX")].(string)),
	)

	for _, record := range event.Records {

		output, err := rk.Process(&types.S3Object{
			Bucket: aws.String(record.S3.Bucket.Name),
			Name:   aws.String(record.S3.Object.Key),
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

		if err := ddb.PutMedia(&database.MediaItem{
			Bucket:     record.S3.Bucket.Name,
			S3Key:      record.S3.Object.Key,
			Faces:      output.Faces,
			Labels:     output.Labels,
			Moderation: output.Moderation,
			Text:       output.Text,
		}); err != nil {
			log.WithFields(logrus.Fields{
				"error": err,
			}).Error("error saving media item")
			return err
		}
	}
	return nil
}
