package rekognition

import (
	"context"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/rekognition"
	"github.com/aws/aws-sdk-go-v2/service/rekognition/types"
	"github.com/sirupsen/logrus"
)

type Detection struct {
	Faces      []types.FaceDetail
	Labels     []types.Label
	Moderation []types.ModerationLabel
	Text       []types.TextDetection
}

type Option func(config *Config)

// Configuration structure.
type Config struct {
	region string
	log    *logrus.Logger
	svc    *rekognition.Client
}

func NewImageProcessor(opts ...func(*Config)) *Config {
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
	cfg.svc = rekognition.NewFromConfig(c)
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

func (config *Config) Process(s3Obj *types.S3Object) (*Detection, error) {
	config.log.WithFields(logrus.Fields{
		"s3Obj": s3Obj,
	}).Info("processing image")
	errors := make(chan error)
	facesChannel := make(chan []types.FaceDetail)
	labelsChannel := make(chan []types.Label)
	moderationChannel := make(chan []types.ModerationLabel)
	textChannel := make(chan []types.TextDetection)

	go func() {
		if faces, err := config.svc.DetectFaces(context.TODO(), &rekognition.DetectFacesInput{
			Image: &types.Image{
				S3Object: s3Obj,
			},
		}); err != nil {
			errors <- err
		} else {
			facesChannel <- faces.FaceDetails
		}
	}()

	go func() {
		if labels, err := config.svc.DetectLabels(context.TODO(), &rekognition.DetectLabelsInput{
			Image: &types.Image{
				S3Object: s3Obj,
			},
		}); err != nil {
			errors <- err
		} else {
			labelsChannel <- labels.Labels
		}
	}()

	go func() {
		if moderation, err := config.svc.DetectModerationLabels(context.TODO(), &rekognition.DetectModerationLabelsInput{
			Image: &types.Image{
				S3Object: s3Obj,
			},
		}); err != nil {
			errors <- err
		} else {
			moderationChannel <- moderation.ModerationLabels
		}
	}()

	go func() {
		if text, err := config.svc.DetectText(context.TODO(), &rekognition.DetectTextInput{
			Image: &types.Image{
				S3Object: s3Obj,
			},
		}); err != nil {
			errors <- err
		} else {
			textChannel <- text.TextDetections
		}
	}()

	detection := &Detection{}

	select {
	case err := <-errors:
		return nil, err
	case faces := <-facesChannel:
		detection.Faces = faces
	case labels := <-labelsChannel:
		detection.Labels = labels
	case moderation := <-moderationChannel:
		detection.Moderation = moderation
	case text := <-textChannel:
		detection.Text = text
	}

	return detection, nil
}
