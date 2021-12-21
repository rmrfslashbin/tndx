package main

import (
	"context"
	"errors"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/rmrfslashbin/tndx/pkg/database"
	"github.com/rmrfslashbin/tndx/pkg/queue"
	"github.com/rmrfslashbin/tndx/pkg/ssmparams"
	"github.com/sirupsen/logrus"
)

var (
	aws_region string
	log        *logrus.Logger
	db         *database.DDBDriver
)

type Message struct {
	RunnerName       string `json:"runner_name"`
	Function         string `json:"function"`
	DDBParamsTable   string `json:"ddb_params_table"`
	DDBRunnerTable   string `json:"ddb_runner_table"`
	SQSRunnerURL     string `json:"sqs_runner_url"`
	S3Bucket         string `json:"s3_bucket"`
	TwitterAPIKey    string `json:"twitter_api_key"`
	TwitterAPISecret string `json:"twitter_api_secret"`
}

func init() {
	log = logrus.New()
	log.SetLevel(logrus.InfoLevel)
	log.SetFormatter(&logrus.JSONFormatter{})
	aws_region = os.Getenv("AWS_REGION")
}

func main() {
	// Catch errors
	var err error
	defer func() {
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"error": err,
			}).Fatal("main crashed")
		}
	}()
	lambda.Start(handler)
}

func handler(ctx context.Context, message Message) error {

	if message.RunnerName == "" {
		return errors.New("runner name is required")
	}
	if message.Function == "" {
		return errors.New("function is required")
	}
	if message.DDBParamsTable == "" {
		return errors.New("ddb params table is required")
	}
	if message.DDBRunnerTable == "" {
		return errors.New("ddb runner table is required")
	}
	if message.SQSRunnerURL == "" {
		return errors.New("sqs runner url is required")
	}
	if message.S3Bucket == "" {
		return errors.New("s3 bucket is required")
	}
	if message.TwitterAPIKey == "" {
		return errors.New("twitter api key is required")
	}
	if message.TwitterAPISecret == "" {
		return errors.New("twitter api secret is required")
	}

	params := ssmparams.NewSSMParams(
		ssmparams.SetRegion(aws_region),
		ssmparams.SetLogger(log),
	)

	outputs, err := params.GetParams([]string{
		message.DDBParamsTable,
		message.DDBRunnerTable,
		message.SQSRunnerURL,
		message.S3Bucket,
		message.TwitterAPIKey,
		message.TwitterAPISecret,
	})

	if err != nil {
		log.WithFields(logrus.Fields{
			"action": "getParams",
			"error":  err.Error(),
		}).Error("error getting parameters.")
		return err
	}

	if len(outputs.InvalidParameters) > 0 {
		log.WithFields(logrus.Fields{
			"invalid_parameters": outputs.InvalidParameters,
		}).Error("invalid parameters")
		return errors.New("invalid parameters")
	}

	db = database.NewDDB(
		database.SetDDBLogger(log),
		database.SetDDBTable(outputs.Params[message.DDBParamsTable].(string)),
		database.SetDDBRunnerTable(outputs.Params[message.DDBRunnerTable].(string)),
	)

	q := queue.NewSQS(
		queue.SetLogger(log),
		queue.SetSQSURL(outputs.Params[message.SQSRunnerURL].(string)),
	)

	users, err := db.GetRunnerUsers(&database.RunnerItem{RunnerName: message.RunnerName})
	if err != nil {
		log.WithFields(logrus.Fields{
			"action": "getRunnerUsers",
			"error":  err.Error(),
		}).Error("error getting runner users.")
		return err
	}

	bootstrap := &queue.Bootstrap{
		S3_bucket:          outputs.Params[message.S3Bucket].(string),
		S3_region:          aws_region,
		DDB_table:          outputs.Params[message.DDBParamsTable].(string),
		DDB_region:         aws_region,
		Twitter_api_key:    outputs.Params[message.TwitterAPIKey].(string),
		Twitter_api_secret: outputs.Params[message.TwitterAPISecret].(string),
		EntiryURL:          "void", // Not needed for this stage
		UserID:             0,      // Not needed for this stage
		TweetId:            "void", // Not needed for this stage
	}

	switch message.Function {
	case "favorites":
		bootstrap.Function = "favorites"
		for _, user := range users {
			if dataabse.Has(user.Flags, database.F_favorites) {
				bootstrap.UserID = user.UserID
				if err := q.SendRunnerMessage(bootstrap); err != nil {
					log.WithFields(logrus.Fields{
						"action": "sendRunnerMessage",
						"error":  err.Error(),
					}).Error("error sending runner message.")
					return err
				}
				log.WithFields(logrus.Fields{
					"bootstrap": bootstrap,
				}).Info("favorites sent.")
			}
		}

	case "followers":
		bootstrap.Function = "followers"
		for _, user := range users {
			if Has(user.Flags, database.F_followers) {
				bootstrap.UserID = user.UserID
				if err := q.SendRunnerMessage(bootstrap); err != nil {
					log.WithFields(logrus.Fields{
						"action": "sendRunnerMessage",
						"error":  err.Error(),
					}).Error("error sending runner message.")
					return err
				}
				log.WithFields(logrus.Fields{
					"bootstrap": bootstrap,
				}).Info("followers sent.")
			}
		}
	case "friends":
		bootstrap.Function = "friends"
		for _, user := range users {
			if Has(user.Flags, database.F_friends) {
				bootstrap.UserID = user.UserID
				if err := q.SendRunnerMessage(bootstrap); err != nil {
					log.WithFields(logrus.Fields{
						"action": "sendRunnerMessage",
						"error":  err.Error(),
					}).Error("error sending runner message.")
					return err
				}
				log.WithFields(logrus.Fields{
					"bootstrap": bootstrap,
				}).Info("friends sent.")
			}
		}
	case "timeline":
		bootstrap.Function = "timeline"
		for _, user := range users {
			if Has(user.Flags, database.F_timeline) {
				bootstrap.UserID = user.UserID
				if err := q.SendRunnerMessage(bootstrap); err != nil {
					log.WithFields(logrus.Fields{
						"action": "sendRunnerMessage",
						"error":  err.Error(),
					}).Error("error sending runner message.")
					return err
				}
				log.WithFields(logrus.Fields{
					"bootstrap": bootstrap,
				}).Info("timeline sent.")
			}
		}
	case "user":
		bootstrap.Function = "user"
		for _, user := range users {
			if Has(user.Flags, database.F_user) {
				bootstrap.UserID = user.UserID
				if err := q.SendRunnerMessage(bootstrap); err != nil {
					log.WithFields(logrus.Fields{
						"action": "sendRunnerMessage",
						"error":  err.Error(),
					}).Error("error sending runner message.")
					return err
				}
				log.WithFields(logrus.Fields{
					"bootstrap": bootstrap,
				}).Info("user sent.")
			}
		}
	default:
		logrus.WithFields(logrus.Fields{
			"function": message.Function,
		}).Error("invalid function; should be one of user, friend, followers, favorites, timeline")
		return errors.New("invalid function; should be one of user, friend, followers, favorites, timeline")
	}

	return nil
}

func getParams(paramNames []*string) (map[string]interface{}, error) {
	s := ssm.New(session.Must(session.NewSession()))
	// Create a SSM client with additional configuration
	//svc := ssm.New(mySession, aws.NewConfig().WithRegion("us-west-2"))

	ret, err := s.GetParameters(&ssm.GetParametersInput{
		Names: paramNames,
	})
	if err != nil {
		log.WithFields(logrus.Fields{
			"action": "ssmparams::GetParameters",
			"error":  err.Error(),
		}).Error("error getting parameters.")
		return nil, err
	}
	output := make(map[string]interface{})

	for _, v := range ret.Parameters {
		output[*v.Name] = *v.Value
	}
	return output, nil

}
