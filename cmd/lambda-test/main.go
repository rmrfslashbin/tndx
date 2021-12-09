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

/*
{
	"function": "user",
	"userid": 242374336,
	"entity_queue": "/tndx/sqs/entities/url",
	"s3_bucket": "/tndx/s3/bucket",
	"s3_region": "/tndx/s3/region",
	"ddb_table": "/tndx/ddb/table",
	"ddb_region": "/tndx/ddb/region",
	"twitter_api_key": "/tndx/twitter/tndx/api/key",
	"twitter_api_secret": "/tndx/twitter/tndx/api/secret"
}
*/

func main() {
	q := queue.NewSQS(
		queue.SetLogger(log),
		queue.SetSQSURL("https://sqs.us-east-1.amazonaws.com/150319663043/tndx-runner"),
	)
	if err := q.SendRunnerMessage(&queue.Bootstrap{
		Function:           "user",
		UserID:             242374336,
		Loglevel:           "info",
		Entity_queue:       "/tndx/sqs/entities/url",
		S3_bucket:          "/tndx/s3/bucket",
		S3_region:          "/tndx/s3/region",
		DDB_table:          "/tndx/ddb/table",
		DDB_region:         "/tndx/ddb/region",
		Twitter_api_key:    "/tndx/twitter/tndx/api/key",
		Twitter_api_secret: "/tndx/twitter/tndx/api/secret",
	}); err != nil {
		log.Error(err)
	}
	log.Info("Sent message")
}
