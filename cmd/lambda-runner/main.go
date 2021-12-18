package main

import (
	"context"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/davecgh/go-spew/spew"
	"github.com/sirupsen/logrus"
)

type Response struct {
	Message string `json:"message"`
}

var log *logrus.Logger

type Message struct {
	RunnerName string `json:"runnername"`
	Function   string `json:"function"`
}

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

func handler(ctx context.Context, message Message) error {

	spew.Dump(message)

	return nil
}
