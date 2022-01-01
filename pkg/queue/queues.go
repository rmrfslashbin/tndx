package queue

import (
	"context"
	"encoding/json"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/sirupsen/logrus"
)

type Bootstrap struct {
	Function         string `json:"function"`
	DDBTablePrefix   string `json:"ddb_table_prefix"`
	DeliveryStream   string `json:"delivery_stream"`
	SQSRunnerURL     string `json:"sqs_runner_url"`
	S3Bucket         string `json:"s3_bucket"`
	TwitterAPIKey    string `json:"twitter_api_key"`
	TwitterAPISecret string `json:"twitter_api_secret"`
}

type ProcessorMessage struct {
	UserID    int64  `json:"user_id"`
	TweetID   string `json:"tweet_id"`
	EntityURL string `json:"entity_url"`
}

type SendMessage struct {
	Bootstrap *Bootstrap        `json:"bootstrap"`
	Message   *ProcessorMessage `json:"message"`
}

type Option func(config *Config)

// Configuration structure.
type Config struct {
	sqsQueueURL string
	region      string
	profile     string
	log         *logrus.Logger
	sqs         *sqs.Client
}

func NewSQS(opts ...func(*Config)) *Config {
	cfg := &Config{}

	// apply the list of options to Config
	for _, opt := range opts {
		opt(cfg)
	}

	if cfg.region == "" {
		cfg.region = os.Getenv("AWS_REGION")
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
	svc := sqs.NewFromConfig(c)
	cfg.sqs = svc

	return cfg
}

func SetProfile(profile string) Option {
	return func(config *Config) {
		config.profile = profile
	}
}

func SetSQSURL(sqsQueueURL string) Option {
	return func(config *Config) {
		config.sqsQueueURL = sqsQueueURL
	}
}

func SetRegion(region string) Option {
	return func(config *Config) {
		config.region = region
	}
}

func SetLogger(log *logrus.Logger) Option {
	return func(config *Config) {
		config.log = log
	}
}

func (config *Config) SendRunnerMessage(params *SendMessage) error {
	body, err := json.Marshal(params.Message)
	if err != nil {
		return err
	}
	message := &sqs.SendMessageInput{
		QueueUrl: aws.String(config.sqsQueueURL),
		MessageAttributes: map[string]types.MessageAttributeValue{
			"function":           {DataType: aws.String("String"), StringValue: aws.String(params.Bootstrap.Function)},
			"ddb_table_prefix":   {DataType: aws.String("String"), StringValue: aws.String(params.Bootstrap.DDBTablePrefix)},
			"delivery_stream":    {DataType: aws.String("String"), StringValue: aws.String(params.Bootstrap.DeliveryStream)},
			"sqs_runner_url":     {DataType: aws.String("String"), StringValue: aws.String(params.Bootstrap.SQSRunnerURL)},
			"s3_bucket":          {DataType: aws.String("String"), StringValue: aws.String(params.Bootstrap.S3Bucket)},
			"twitter_api_key":    {DataType: aws.String("String"), StringValue: aws.String(params.Bootstrap.TwitterAPIKey)},
			"twitter_api_secret": {DataType: aws.String("String"), StringValue: aws.String(params.Bootstrap.TwitterAPISecret)},
		},
		MessageBody: aws.String(string(body)),
	}
	/*
		config.log.WithFields(logrus.Fields{
			"message": message,
		}).Info("Sending message to SQS")
	*/
	_, err = config.sqs.SendMessage(context.TODO(), message)
	return err
}

func (config *Config) GetAttribs() (map[string]string, error) {
	message := &sqs.GetQueueAttributesInput{
		QueueUrl: aws.String(config.sqsQueueURL),
		AttributeNames: []types.QueueAttributeName{
			"All",
		},
	}
	if ret, err := config.sqs.GetQueueAttributes(context.TODO(), message); err != nil {
		return nil, err
	} else {
		return ret.Attributes, nil
	}
}
