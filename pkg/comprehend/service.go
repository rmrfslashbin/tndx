package comprehend

import (
	"context"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/comprehend"
	"github.com/aws/aws-sdk-go-v2/service/comprehend/types"
	"github.com/sirupsen/logrus"
)

type Option func(config *Config)

// Configuration structure.
type Config struct {
	region  string
	profile string
	log     *logrus.Logger
	svc     *comprehend.Client
}

func New(opts ...func(*Config)) *Config {
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

	cfg.svc = comprehend.NewFromConfig(c)

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

func (config *Config) DetectEntities(textList *[]string) (*comprehend.BatchDetectEntitiesOutput, error) {
	return config.svc.BatchDetectEntities(
		context.TODO(),
		&comprehend.BatchDetectEntitiesInput{
			LanguageCode: types.LanguageCodeEn,
			TextList:     *textList,
		})
}

func (config *Config) Sentiment(textList *[]string) (*comprehend.BatchDetectSentimentOutput, error) {
	return config.svc.BatchDetectSentiment(
		context.TODO(),
		&comprehend.BatchDetectSentimentInput{
			LanguageCode: types.LanguageCodeEn,
			TextList:     *textList,
		})
}
