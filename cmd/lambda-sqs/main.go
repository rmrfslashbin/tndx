package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/sirupsen/logrus"
)

type Response struct {
	Message string `json:"message"`
}

var log *logrus.Logger

type Messages []sqs.Message

func init() {
	log = logrus.New()
	log.SetLevel(logrus.InfoLevel)
	log.SetFormatter(&logrus.JSONFormatter{})
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

func handler(ctx context.Context, sqsEvent events.SQSEvent) error {

	for _, message := range sqsEvent.Records {
		for i, mattr := range message.MessageAttributes {
			fmt.Printf("%s: %s\n", i, *mattr.StringValue)
		}
	}

	return nil
}
