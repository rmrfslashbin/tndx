package events

import (
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/eventbridge"
	"github.com/sirupsen/logrus"
)

type Option func(config *Config)

// Configuration structure.
type Config struct {
	region string
	log    *logrus.Logger
	svc    *eventbridge.EventBridge
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

	cfg.svc = eventbridge.New(session.Must(session.NewSession()), aws.NewConfig().WithRegion(cfg.region))

	return cfg
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
	return config.svc.ListRules(&eventbridge.ListRulesInput{
		NamePrefix: aws.String("tndx-"),
	})
}

func (config *Config) Disable(ruleName *string) error {
	_, err := config.svc.DisableRule(&eventbridge.DisableRuleInput{
		Name: ruleName,
	})
	return err
}

func (config *Config) Enable(ruleName *string) error {
	_, err := config.svc.EnableRule(&eventbridge.EnableRuleInput{
		Name: ruleName,
	})
	return err
}
