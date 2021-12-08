package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/davecgh/go-spew/spew"
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
	lambda.Start(HandleLambdaEvent)
}

func HandleLambdaEvent(event Messages) (Response, error) {
	spew.Dump(event)

	return Response{
		Message: "All done!",
	}, nil
}
