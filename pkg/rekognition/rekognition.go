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
	region  string
	profile string
	log     *logrus.Logger
	svc     *rekognition.Client
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
		if cfg.profile != "" {
			o.SharedConfigProfile = cfg.profile
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
	cfg.svc = rekognition.NewFromConfig(c)
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

func (config *Config) Process(s3Obj *types.S3Object) (*Detection, error) {
	type Faces struct {
		faces []types.FaceDetail
		err   error
	}
	type Labels struct {
		labels []types.Label
		err    error
	}

	type Moderation struct {
		moderation []types.ModerationLabel
		err        error
	}

	type Text struct {
		text []types.TextDetection
		err  error
	}

	config.log.WithFields(logrus.Fields{
		"s3Obj": s3Obj,
	}).Info("processing image")

	facesChannel := make(chan Faces)
	labelsChannel := make(chan Labels)
	moderationChannel := make(chan Moderation)
	textChannel := make(chan Text)

	go func() {
		if faces, err := config.svc.DetectFaces(context.TODO(), &rekognition.DetectFacesInput{
			Image: &types.Image{
				S3Object: s3Obj,
			},
		}); err != nil {
			facesChannel <- Faces{err: err}
		} else {
			facesChannel <- Faces{faces: faces.FaceDetails}
		}
	}()

	go func() {
		if labels, err := config.svc.DetectLabels(context.TODO(), &rekognition.DetectLabelsInput{
			Image: &types.Image{
				S3Object: s3Obj,
			},
		}); err != nil {
			labelsChannel <- Labels{err: err}
		} else {
			labelsChannel <- Labels{labels: labels.Labels}
		}
	}()

	go func() {
		if moderation, err := config.svc.DetectModerationLabels(context.TODO(), &rekognition.DetectModerationLabelsInput{
			Image: &types.Image{
				S3Object: s3Obj,
			},
		}); err != nil {
			moderationChannel <- Moderation{err: err}
		} else {
			moderationChannel <- Moderation{moderation: moderation.ModerationLabels}
		}
	}()

	go func() {
		if text, err := config.svc.DetectText(context.TODO(), &rekognition.DetectTextInput{
			Image: &types.Image{
				S3Object: s3Obj,
			},
		}); err != nil {
			textChannel <- Text{err: err}
		} else {
			textChannel <- Text{text: text.TextDetections}
		}
	}()

	faces, labels, moderation, text := <-facesChannel, <-labelsChannel, <-moderationChannel, <-textChannel
	if faces.err != nil {
		config.log.WithFields(logrus.Fields{
			"error": faces.err,
		}).Error("error processing faces")
		return nil, faces.err
	}
	if labels.err != nil {
		config.log.WithFields(logrus.Fields{
			"error": labels.err,
		}).Error("error processing labels")
		return nil, labels.err
	}
	if moderation.err != nil {
		config.log.WithFields(logrus.Fields{
			"error": moderation.err,
		}).Error("error processing moderation")
		return nil, moderation.err
	}
	if text.err != nil {
		config.log.WithFields(logrus.Fields{
			"error": text.err,
		}).Error("error processing text")
		return nil, text.err
	}

	return &Detection{
		Faces:      faces.faces,
		Labels:     labels.labels,
		Moderation: moderation.moderation,
		Text:       text.text,
	}, nil
}
