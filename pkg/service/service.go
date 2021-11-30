package service

import (
	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"github.com/sirupsen/logrus"
)

// User query params.
type QueryParams struct {
	ScreenName          string
	UserID              int64
	Count               int
	SinceID             int64
	MaxID               int64
	Cursor              int64
	SkipStatus          *bool
	IncludeUserEntities *bool
}

type Option func(config *Config)

// Configuration structure.
type Config struct {
	accessSecret   string
	accessToken    string
	consumerKey    string
	consumerSecret string
	log            *logrus.Logger
	client         *twitter.Client
}

// New is a factory function for creating a new Config
func New(opts ...func(*Config)) *Config {
	config := &Config{}

	// apply the list of options to Config
	for _, opt := range opts {
		opt(config)
	}

	oauthConfig := oauth1.NewConfig(config.consumerKey, config.consumerSecret)
	oauthToken := oauth1.NewToken(config.accessToken, config.accessSecret)
	// http.Client will automatically authorize Requests
	httpClient := oauthConfig.Client(oauth1.NoContext, oauthToken)

	// Twitter client
	config.client = twitter.NewClient(httpClient)
	return config
}

func SetAccessSecret(accessSecret string) Option {
	return func(config *Config) {
		config.accessSecret = accessSecret
	}
}

func SetAccessToken(accessToken string) Option {
	return func(config *Config) {
		config.accessToken = accessToken
	}
}

func SetConsumerKey(consumerKey string) Option {
	return func(c *Config) {
		c.consumerKey = consumerKey
	}
}

func SetConsumerSecret(consumerSecret string) Option {
	return func(config *Config) {
		config.consumerSecret = consumerSecret
	}
}

func SetLogger(log *logrus.Logger) Option {
	return func(config *Config) {
		config.log = log
	}
}
