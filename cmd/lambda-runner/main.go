package main

import (
	"context"
	"errors"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/davecgh/go-spew/spew"
	"github.com/rmrfslashbin/tndx/pkg/database"
	"github.com/rmrfslashbin/tndx/pkg/queue"
	"github.com/sirupsen/logrus"
)

var (
	log *logrus.Logger
	db  *database.DDBDriver
)

type Message struct {
	RunnerName       string `json:"runner_name"`
	Loglevel         string `json:"loglevel"`
	Function         string `json:"function"`
	DDBRegion        string `json:"ddb_region"`
	DDBTable         string `json:"ddb_table"`
	DDBRunnerTable   string `json:"ddb_runner_table"`
	SQSRunnerURL     string `json:"sqs_runner_url"`
	SQSEntityURL     string `json:"sqs_entity_url"`
	S3Bucket         string `json:"s3_bucket"`
	S3Region         string `json:"s3_region"`
	TwitterAPIKey    string `json:"twitter_api_key"`
	TwitterAPISecret string `json:"twitter_api_secret"`
}

func init() {
	log = logrus.New()
	log.SetLevel(logrus.InfoLevel)
	log.SetFormatter(&logrus.JSONFormatter{})
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
	outputs, err := getParams([]*string{
		&message.DDBRegion,
		&message.DDBRunnerTable,
		&message.DDBTable,
		&message.SQSRunnerURL,
		&message.SQSEntityURL,
		&message.S3Bucket,
		&message.S3Region,
		&message.TwitterAPIKey,
		&message.TwitterAPISecret,
	})
	if err != nil {
		log.WithFields(logrus.Fields{
			"action": "getParams",
			"error":  err.Error(),
		}).Error("error getting parameters.")
		return err
	}

	db = database.NewDDB(
		database.SetDDBLogger(log),
		database.SetDDBTable(outputs[message.DDBRegion].(string)),
		database.SetDDBRunnerTable(outputs[message.DDBRunnerTable].(string)),
	)

	q := queue.NewSQS(
		queue.SetLogger(log),
		queue.SetSQSURL(outputs[message.SQSRunnerURL].(string)),
	)

	users, err := db.GetRunnerUsers(&database.RunnerUsersInput{RunnerName: message.RunnerName})
	if err != nil {
		log.WithFields(logrus.Fields{
			"action": "getRunnerUsers",
			"error":  err.Error(),
		}).Error("error getting runner users.")
		return err
	}

	bootstrap := &queue.Bootstrap{
		Loglevel:           message.Loglevel,
		Entity_queue:       outputs[message.SQSEntityURL].(string),
		S3_bucket:          outputs[message.S3Bucket].(string),
		S3_region:          outputs[message.S3Region].(string),
		DDB_table:          outputs[message.DDBTable].(string),
		DDB_region:         outputs[message.DDBRegion].(string),
		Twitter_api_key:    outputs[message.TwitterAPIKey].(string),
		Twitter_api_secret: outputs[message.TwitterAPISecret].(string),
	}
	spew.Dump(bootstrap)

	switch message.Function {
	case "favorites":
		bootstrap.Function = "favorites"
		for _, user := range users {
			if Has(user.Flags, database.F_favorites) {
				bootstrap.UserID = user.UserID
				if err := q.SendRunnerMessage(bootstrap); err != nil {
					log.WithFields(logrus.Fields{
						"action": "sendRunnerMessage",
						"error":  err.Error(),
					}).Error("error sending runner message.")
					return err
				}
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

func Set(b, flag database.Bits) database.Bits    { return b | flag }
func Clear(b, flag database.Bits) database.Bits  { return b &^ flag }
func Toggle(b, flag database.Bits) database.Bits { return b ^ flag }
func Has(b, flag database.Bits) bool             { return b&flag != 0 }
