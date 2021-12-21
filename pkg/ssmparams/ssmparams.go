package ssmparams

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/sirupsen/logrus"
)

type SSMParamsOutput struct {
	Params            map[string]interface{}
	InvalidParameters []string
}

type SSMParamsOption func(config *SSMParamsConfig)

type SSMParamsConfig struct {
	log     *logrus.Logger
	region  string
	profile string
	ssm     *ssm.Client
}

func NewSSMParams(opts ...func(*SSMParamsConfig)) *SSMParamsConfig {
	cfg := &SSMParamsConfig{}

	// apply the list of options to Config
	for _, opt := range opts {
		opt(cfg)
	}

	if cfg.region == "" {
		log.Fatal("region is required")
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
	svc := ssm.NewFromConfig(c)
	cfg.ssm = svc

	return cfg
}

func SetLogger(logger *logrus.Logger) SSMParamsOption {
	return func(config *SSMParamsConfig) {
		config.log = logger
	}
}

func SetProfile(profile string) SSMParamsOption {
	return func(config *SSMParamsConfig) {
		config.profile = profile
	}
}

func SetRegion(region string) SSMParamsOption {
	return func(config *SSMParamsConfig) {
		config.region = region
	}
}

func (config *SSMParamsConfig) GetParams(paramNames []string) (*SSMParamsOutput, error) {
	params, err := config.ssm.GetParameters(context.TODO(), &ssm.GetParametersInput{
		Names: paramNames,
	})
	if err != nil {
		config.log.WithFields(logrus.Fields{
			"action": "ssmparams::GetParameters",
			"error":  err.Error(),
		}).Error("error getting parameters.")
		return nil, err
	}
	output := make(map[string]interface{}, len(params.Parameters))

	for _, v := range params.Parameters {
		output[*v.Name] = *v.Value
	}
	return &SSMParamsOutput{
		Params:            output,
		InvalidParameters: params.InvalidParameters,
	}, nil
}
