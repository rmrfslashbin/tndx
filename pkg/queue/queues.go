package queue

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/sirupsen/logrus"
)

type SSMParams struct {
	EntityQueue      string `json:"entity_queue"`
	S3Bucket         string `json:"s3_bucket"`
	S3Region         string `json:"s3_region"`
	DDBTable         string `json:"ddb_table"`
	DDBRegion        string `json:"ddb_region"`
	TwitterAPIKey    string `json:"twitter_api_key"`
	TwitterAPISecret string `json:"twitter_api_secret"`
}

type Bootstrap struct {
	Function           string `json:"function"` // user, friends, followers, favorties, timeline, entities
	Loglevel           string `json:"loglevel"` // error, warn, info, debug, trace
	UserID             int64  `json:"userid"`
	Entity_queue       string `json:"entity_queue"`
	S3_bucket          string `json:"s3_bucket"`
	S3_region          string `json:"s3_region"`
	DDB_table          string `json:"ddb_table"`
	DDB_region         string `json:"ddb_region"`
	Twitter_api_key    string `json:"twitter_api_key"`
	Twitter_api_secret string `json:"twitter_api_secret"`
}

type Option func(config *Config)

// Configuration structure.
type Config struct {
	sqsQueueURL string
	s3Bucket    string
	log         *logrus.Logger
	sqs         *sqs.SQS
}

func NewSQS(opts ...func(*Config)) *Config {
	config := &Config{}

	// apply the list of options to Config
	for _, opt := range opts {
		opt(config)
	}

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	// Create DynamoDB client
	svc := sqs.New(sess)
	config.sqs = svc

	return config
}

func SetSQSURL(sqsQueueURL string) Option {
	return func(config *Config) {
		config.sqsQueueURL = sqsQueueURL
	}
}

func SetS3Bucket(s3Bucket string) Option {
	return func(config *Config) {
		config.s3Bucket = s3Bucket
	}
}

func SetLogger(log *logrus.Logger) Option {
	return func(config *Config) {
		config.log = log
	}
}
