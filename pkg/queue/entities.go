package queue

import (
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/rmrfslashbin/tndx/pkg/utils"
)

func (config *Config) SendEntityMessage(tweetId string, url string) error {
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
			"bucket": &sqs.MessageAttributeValue{
				DataType:    aws.String("String"),
				StringValue: aws.String(config.s3Bucket),
			},
		},
		MessageBody: aws.String("This is a test message."),
		QueueUrl:    &config.sqsQueueURL,
	})
	return err
}

func (config *Config) SendRunnerMessage(Bootstrap *Bootstrap) error {
	_, err := config.sqs.SendMessage(&sqs.SendMessageInput{
		DelaySeconds: aws.Int64(10),
		MessageAttributes: map[string]*sqs.MessageAttributeValue{
			"function": &sqs.MessageAttributeValue{
				DataType:    aws.String("String"),
				StringValue: aws.String(Bootstrap.Function),
			},
			"loglevel": &sqs.MessageAttributeValue{
				DataType:    aws.String("String"),
				StringValue: aws.String(Bootstrap.Loglevel),
			},
			"userid": &sqs.MessageAttributeValue{
				DataType:    aws.String("Number"),
				StringValue: aws.String(strconv.FormatInt(Bootstrap.UserID, 10)),
			},
			"entity_queue": &sqs.MessageAttributeValue{
				DataType:    aws.String("String"),
				StringValue: aws.String(Bootstrap.Entity_queue),
			},
			"s3_bucket": &sqs.MessageAttributeValue{
				DataType:    aws.String("String"),
				StringValue: aws.String(Bootstrap.S3_bucket),
			},
			"s3_region": &sqs.MessageAttributeValue{
				DataType:    aws.String("String"),
				StringValue: aws.String(Bootstrap.S3_region),
			},
			"ddb_table": &sqs.MessageAttributeValue{
				DataType:    aws.String("String"),
				StringValue: aws.String(Bootstrap.DDB_table),
			},
			"ddb_region": &sqs.MessageAttributeValue{
				DataType:    aws.String("String"),
				StringValue: aws.String(Bootstrap.DDB_region),
			},
			"twitter_api_key": &sqs.MessageAttributeValue{
				DataType:    aws.String("String"),
				StringValue: aws.String(Bootstrap.Twitter_api_key),
			},
			"twitter_api_secret": &sqs.MessageAttributeValue{
				DataType:    aws.String("String"),
				StringValue: aws.String(Bootstrap.Twitter_api_secret),
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
