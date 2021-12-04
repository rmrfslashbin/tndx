package queue

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/rmrfslashbin/tndx/pkg/utils"
	"github.com/sirupsen/logrus"
)

type Option func(config *Config)

// Configuration structure.
type Config struct {
	sqsQueueURL string
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

func SetLogger(log *logrus.Logger) Option {
	return func(config *Config) {
		config.log = log
	}
}

func (config *Config) SendMessage(tweetId string, url string) error {
	_, err := config.sqs.SendMessage(&sqs.SendMessageInput{
		DelaySeconds: aws.Int64(10),
		MessageAttributes: map[string]*sqs.MessageAttributeValue{
			"tweetId": &sqs.MessageAttributeValue{
				DataType:    aws.String("String"),
				StringValue: aws.String(tweetId),
			},
			"url": &sqs.MessageAttributeValue{
				DataType:    aws.String("String"),
				StringValue: aws.String(url),
			},
		},
		MessageBody: aws.String("This is a test message."),
		QueueUrl:    &config.sqsQueueURL,
	})
	return err
}

func (config *Config) ReceiveMessage() error {
	result, err := config.sqs.ReceiveMessage(&sqs.ReceiveMessageInput{
		AttributeNames: []*string{
			aws.String(sqs.MessageSystemAttributeNameSentTimestamp),
		},
		MessageAttributeNames: []*string{
			aws.String(sqs.QueueAttributeNameAll),
		},
		QueueUrl:            &config.sqsQueueURL,
		MaxNumberOfMessages: aws.Int64(1),
		VisibilityTimeout:   aws.Int64(30),
	})
	if err != nil {
		config.log.Error(err)
		return err
	}
	for _, message := range result.Messages {
		userid := (message.MessageAttributes["tweetId"].StringValue)
		url := message.MessageAttributes["url"].StringValue
		if err := utils.SaveEntities(userid, url); err != nil {
			config.log.Error(err)
			return err
		}
	}
	return nil
}
