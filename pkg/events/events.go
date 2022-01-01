package events

import (
	"context"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/eventbridge"
	"github.com/sirupsen/logrus"
)

type Option func(config *Config)

// Configuration structure.
type Config struct {
	region  string
	profile string
	log     *logrus.Logger
	svc     *eventbridge.Client
}

func NewEvents(opts ...func(*Config)) *Config {
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

	cfg.svc = eventbridge.NewFromConfig(c)

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

func (config *Config) List() (*eventbridge.ListRulesOutput, error) {
	return config.svc.ListRules(context.TODO(), &eventbridge.ListRulesInput{
		NamePrefix: aws.String("tndx-"),
	})

}

func (config *Config) Disable(ruleName *string) error {
	_, err := config.svc.DisableRule(context.TODO(), &eventbridge.DisableRuleInput{
		Name: ruleName,
	})
	return err
}

func (config *Config) Enable(ruleName *string) error {
	_, err := config.svc.EnableRule(context.TODO(), &eventbridge.EnableRuleInput{
		Name: ruleName,
	})
	return err
}
