package main

import (
	"github.com/rmrfslashbin/tndx/pkg/queue"
	"github.com/sirupsen/logrus"
)

var log *logrus.Logger

func init() {
	log = logrus.New()
	log.SetLevel(logrus.InfoLevel)
	log.SetFormatter(&logrus.JSONFormatter{})
}

func main() {
	q := queue.NewSQS(
		queue.SetLogger(log),
		queue.SetSQSURL("https://sqs.us-east-1.amazonaws.com/150319663043/tndx-runner"),
	)
	if err := q.SendRunnerMessage(&queue.Bootstrap{
		Function: "func01",
		UserID:   1234,
		Loglevel: "info",
	}); err != nil {
		log.Error(err)
	}
	log.Info("Sent message")
}
