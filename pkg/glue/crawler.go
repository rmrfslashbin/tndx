package glue

import (
	"context"
	"errors"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/glue"
	"github.com/sirupsen/logrus"
)

type Option func(config *Config)

// Configuration structure.
type Config struct {
	region      string
	log         *logrus.Logger
	crawlerName string
	glue        *glue.Client
}

func NewCrawler(opts ...func(*Config)) *Config {
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
		return nil
	})
	if err != nil {
		panic(err)
	}
	cfg.glue = glue.NewFromConfig(c)
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

func SetCrawlerName(crawlerName string) Option {
	return func(config *Config) {
		config.crawlerName = crawlerName
	}
}

func (config *Config) GetCrawlerData() (*glue.GetCrawlerOutput, error) {
	if config.crawlerName == "" {
		return nil, errors.New("crawler name is not set")
	}

	return config.glue.GetCrawler(context.TODO(), &glue.GetCrawlerInput{
		Name: aws.String(config.crawlerName),
	})
}

func (config *Config) ListCrawlers() ([]string, error) {
	if ret, err := config.glue.ListCrawlers(context.TODO(), &glue.ListCrawlersInput{}); err != nil {
		return nil, err
	} else {
		return ret.CrawlerNames, nil
	}
}

func (config *Config) StartCrawler() error {
	if config.crawlerName == "" {
		return errors.New("crawler name is not set")
	}

	if _, err := config.glue.StartCrawler(context.TODO(), &glue.StartCrawlerInput{
		Name: aws.String(config.crawlerName),
	}); err != nil {
		return err
	}
	return nil
}
