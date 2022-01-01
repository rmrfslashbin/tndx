package kenisis

import (
	"context"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/firehose"
	"github.com/aws/aws-sdk-go-v2/service/firehose/types"
	"github.com/sirupsen/logrus"
)

type Option func(config *Config)

// Configuration structure.
type Config struct {
	region         string
	profile        string
	log            *logrus.Logger
	firehose       *firehose.Client
	deliveryStream *string
}

func NewFirehose(opts ...func(*Config)) *Config {
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
	svc := firehose.NewFromConfig(c)
	cfg.firehose = svc

	return cfg
}

func SetProfile(profile string) Option {
	return func(config *Config) {
		config.profile = profile
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

func SetDeliveryStream(deliveryStream string) Option {
	return func(config *Config) {
		config.deliveryStream = &deliveryStream
	}
}

func (config *Config) PutRecord(data []byte) (*firehose.PutRecordOutput, error) {
	record := &firehose.PutRecordInput{
		DeliveryStreamName: config.deliveryStream,
		Record: &types.Record{
			Data: data,
		},
	}
	return config.firehose.PutRecord(context.TODO(), record)
}
