package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"path"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/rmrfslashbin/tndx/pkg/database"
	"github.com/rmrfslashbin/tndx/pkg/queue"
	"github.com/rmrfslashbin/tndx/pkg/service"
	"github.com/rmrfslashbin/tndx/pkg/storage"
	"github.com/sirupsen/logrus"
)

// service stores drivers and clients
type services struct {
	twitterClient *service.Config
	storage       storage.StorageDriver
	db            database.DatabaseDriver
	queue         *queue.Config
}

type RunnerFunction struct {
	Function         *string `json:"function"`
	DDBRegion        *string `json:"ddb_region"`
	DDBTable         *string `json:"ddb_table"`
	SQSEntityURL     *string `json:"sqs_entity_url"`
	S3Bucket         *string `json:"s3_bucket"`
	S3Region         *string `json:"s3_region"`
	TwitterAPIKey    *string `json:"twitter_api_key"`
	TwitterAPISecret *string `json:"twitter_api_secret"`
	UserID           int64   `json:"userid"`
}

var (
	log *logrus.Logger
	svc services
)

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

func handler(ctx context.Context, sqsEvent events.SQSEvent) error {

	for _, message := range sqsEvent.Records {
		userId, err := strconv.ParseInt(*message.MessageAttributes["userid"].StringValue, 10, 64)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"action": "handler::ParseInt",
				"error":  err.Error(),
			}).Error("error parsing userid")
			return err
		}
		runnerFunction := &RunnerFunction{
			Function:         message.MessageAttributes["function"].StringValue,
			DDBRegion:        message.MessageAttributes["ddb_region"].StringValue,
			DDBTable:         message.MessageAttributes["ddb_table"].StringValue,
			SQSEntityURL:     message.MessageAttributes["sqs_entity_url"].StringValue,
			S3Bucket:         message.MessageAttributes["s3_bucket"].StringValue,
			S3Region:         message.MessageAttributes["s3_region"].StringValue,
			TwitterAPIKey:    message.MessageAttributes["twitter_api_key"].StringValue,
			TwitterAPISecret: message.MessageAttributes["twitter_api_secret"].StringValue,
			UserID:           userId,
		}

		svc.queue = queue.NewSQS(
			queue.SetLogger(log),
			queue.SetSQSURL(*runnerFunction.SQSEntityURL),
			queue.SetS3Bucket(*runnerFunction.S3Bucket),
		)

		svc.db = database.NewDDB(
			database.SetDDBLogger(log),
			database.SetDDBTable(*runnerFunction.DDBTable),
			database.SetDDBRegion(*runnerFunction.DDBRegion),
		)

		svc.storage = storage.NewS3Storage(
			storage.SetS3Bucket(*runnerFunction.S3Bucket),
			storage.SetS3Region(*runnerFunction.S3Region),
		)

		svc.twitterClient = service.New(
			service.SetConsumerKey(*runnerFunction.TwitterAPIKey),
			service.SetConsumerSecret(*runnerFunction.TwitterAPISecret),
			service.SetLogger(log),
		)
		switch *runnerFunction.Function {
		case "favorites":
			if err := favorites(runnerFunction.UserID); err != nil {
				logrus.WithFields(logrus.Fields{
					"function": "favorites",
					"error":    err,
				}).Error("function failed")
				return err
			}

		case "followers":
			if err := followers(runnerFunction.UserID); err != nil {
				logrus.WithFields(logrus.Fields{
					"function": "followers",
					"error":    err,
				}).Error("function failed")
				return err
			}

		case "friends":
			if err := friends(runnerFunction.UserID); err != nil {
				logrus.WithFields(logrus.Fields{
					"function": "friends",
					"error":    err,
				}).Error("function failed")
				return err
			}

		case "timeline":
			if err := timeline(runnerFunction.UserID); err != nil {
				logrus.WithFields(logrus.Fields{
					"function": "timeline",
					"error":    err,
				}).Error("function failed")
				return err
			}

		case "user":
			if err := user(runnerFunction.UserID); err != nil {
				logrus.WithFields(logrus.Fields{
					"function": "user",
					"error":    err,
				}).Error("function failed")
				return err
			}

		case "entities":
			logrus.WithFields(logrus.Fields{
				"function": "entities",
			}).Info("Entities could be processed here.")

		default:
			logrus.WithFields(logrus.Fields{
				"function": *runnerFunction.Function,
			}).Error("invalid function; should be one of user, friend, followers, favorites, timeline, entities")
			return errors.New("invalid function; should be one of user, friend, followers, favorites, timeline, entities")
		}
	}

	return nil
}

func favorites(userid int64) error {
	favConfig, err := svc.db.GetFavoritesConfig(userid)
	if err != nil {
		log.WithFields(logrus.Fields{
			"action": "favorites::GetFavoritesConfig",
			"error":  err.Error(),
		}).Error("error getting favorites config")
		return err
	}

	log.WithFields(logrus.Fields{
		"action":  "favorites::Setup",
		"userid":  userid,
		"sinceid": favConfig.MaxID,
	}).Info("setting up favorites")

	tweets, resp, err := svc.twitterClient.GetUserFavorites(
		&service.QueryParams{
			Count:   200,
			SinceID: favConfig.MaxID,
			UserID:  userid,
		},
	)
	if err != nil {
		log.WithFields(logrus.Fields{
			"action":   "favorites",
			"response": resp.Status,
		}).Error("error getting user's favorites")
		return err
	}

	// upperID and lowerID are used to keep track of the max and min tweet IDs
	var upperID int64
	var lowerID int64

	// Loop through all the tweets.
	for t := range tweets {
		if data, err := json.Marshal(tweets[t]); err == nil {
			if err := svc.storage.Put(path.Join("favorites", strconv.FormatInt(userid, 10), tweets[t].IDStr+".json"), data); err != nil {
				logrus.WithFields(logrus.Fields{
					"action":  "favorites::Put",
					"error":   err.Error(),
					"userid":  userid,
					"tweetId": tweets[t].ID,
				}).Error("error putting favorites")
				return err
			}
		}

		// Loop through all the media entities
		for m := range tweets[t].Entities.Media {
			var url string
			if tweets[t].Entities.Media[m].MediaURLHttps != "" {
				url = tweets[t].Entities.Media[m].MediaURLHttps
			} else if tweets[t].Entities.Media[m].MediaURL != "" {
				url = tweets[t].Entities.Media[m].MediaURL
			}
			if url != "" {
				if err := svc.queue.SendEntityMessage(tweets[t].IDStr, url); err != nil {
					logrus.WithFields(logrus.Fields{
						"action":  "favorites::queue::SendMessage",
						"error":   err.Error(),
						"userid":  userid,
						"tweetId": tweets[t].ID,
					}).Error("error sending message to queue")
					fmt.Printf("Queued: %s\n", url)
				}
			}
		}

		// Calculate the max and min tweet IDs.
		if tweets[t].ID > upperID {
			upperID = tweets[t].ID
		}
		if tweets[t].ID < lowerID {
			lowerID = tweets[t].ID
		}
		if lowerID == 0 {
			lowerID = upperID
		}
	}

	if upperID > 0 {
		if err := svc.db.PutFavoritesConfig(
			&database.TweetConfigQuery{
				UserID:  userid,
				SinceID: lowerID,
				MaxID:   upperID,
			},
		); err != nil {
			logrus.WithFields(logrus.Fields{
				"action":       "favorites::PutFavoritesConfig",
				"error":        err.Error(),
				"userid":       userid,
				"MaxUpperID":   upperID,
				"SinceLowerID": lowerID,
			}).Error("error putting favorites config")
			return err
		}
	}

	logrus.WithFields(logrus.Fields{
		"action":  "favorites::Done!",
		"userid":  userid,
		"upperID": upperID,
		"lowerID": lowerID,
		"count":   len(tweets),
	}).Info("finished getting favorites")

	return nil
}

func followers(userid int64) error {
	followersConfig, err := svc.db.GetFollowersConfig(userid)
	if err != nil {
		log.WithFields(logrus.Fields{
			"action": "followers::GetFollowersConfig",
			"error":  err.Error(),
		}).Error("error getting follower config")
		return err
	}

	log.WithFields(logrus.Fields{
		"action": "followers::Setup",
		"userid": userid,
		"cursor": followersConfig.NextCursor,
	}).Debug("setting up followers")

	followers, resp, err := svc.twitterClient.GetUserFollowers(
		&service.QueryParams{
			Count:  200,
			UserID: userid,
			Cursor: followersConfig.NextCursor,
		},
	)
	if err != nil {
		log.WithFields(logrus.Fields{
			"action":   "followers::GetUserFollowers",
			"response": resp.Status,
		}).Error("error getting user's followers")
		return err
	}

	if err := svc.db.PutFollowersConfig(
		&database.CursoredTweetConfigQuery{
			UserID:         userid,
			NextCursor:     followers.NextCursor,
			PreviousCursor: followers.PreviousCursor,
		},
	); err != nil {
		logrus.WithFields(logrus.Fields{
			"action":         "followers::PutFollowersConfig",
			"error":          err.Error(),
			"userid":         userid,
			"nextCursor":     followers.NextCursor,
			"previousCursor": followers.PreviousCursor,
		}).Error("error putting followers config")
		return err
	}

	// Save the users.
	for f := range followers.Users {
		if data, err := json.Marshal(followers.Users[f]); err == nil {
			if err := svc.storage.Put(path.Join("followers", strconv.FormatInt(userid, 10), followers.Users[f].IDStr+".json"), data); err != nil {
				logrus.WithFields(logrus.Fields{
					"action":     "followers::Put",
					"error":      err.Error(),
					"userid":     userid,
					"followerId": followers.Users[f].ID,
				}).Error("error putting followers")
				return err
			}
		}
	}
	logrus.WithFields(logrus.Fields{
		"action":         "followers::Done!",
		"userid":         userid,
		"nextCursor":     followers.NextCursor,
		"previousCursor": followers.PreviousCursor,
		"count":          len(followers.Users),
	}).Info("finished getting followers")

	return nil
}

func friends(userid int64) error {
	friendsConfig, err := svc.db.GetFriendsConfig(userid)
	if err != nil {
		log.WithFields(logrus.Fields{
			"action": "friends::GetFriendsConfig",
			"error":  err.Error(),
		}).Error("error getting friends config")
		return err
	}
	log.WithFields(logrus.Fields{
		"action":     "friends::GetFriendsConfig",
		"userid":     userid,
		"nextCursor": friendsConfig.NextCursor,
	}).Debug("got friends config")

	log.WithFields(logrus.Fields{
		"action": "friends::Setup",
		"userid": userid,
		"cursor": friendsConfig.NextCursor,
	}).Debug("setting up friends")

	friends, resp, err := svc.twitterClient.GetUserFriends(
		&service.QueryParams{
			Count:  200,
			UserID: userid,
			Cursor: friendsConfig.NextCursor,
		},
	)
	if err != nil {
		log.WithFields(logrus.Fields{
			"action":   "friends::GetUserFriends",
			"response": resp.Status,
		}).Error("error getting user's friends")
		return err
	}

	if err := svc.db.PutFriendsConfig(
		&database.CursoredTweetConfigQuery{
			UserID:     userid,
			NextCursor: friends.NextCursor,
		},
	); err != nil {
		logrus.WithFields(logrus.Fields{
			"action":         "friends::PutFriendsConfig",
			"error":          err.Error(),
			"userid":         userid,
			"nextCursor":     friends.NextCursor,
			"previousCursor": friends.PreviousCursor,
		}).Error("error putting friends config")
		return err
	}

	// Save the users.
	for f := range friends.Users {
		if data, err := json.Marshal(friends.Users[f]); err == nil {
			if err := svc.storage.Put(path.Join("friends", strconv.FormatInt(userid, 10), friends.Users[f].IDStr+".json"), data); err != nil {
				logrus.WithFields(logrus.Fields{
					"action":   "friends::Put",
					"error":    err.Error(),
					"userid":   userid,
					"friendId": friends.Users[f].ID,
				}).Error("error putting friends")
				return err
			}
		}
	}
	logrus.WithFields(logrus.Fields{
		"action":         "friends::Done!",
		"userid":         userid,
		"nextCursor":     friends.NextCursor,
		"previousCursor": friends.PreviousCursor,
		"count":          len(friends.Users),
	}).Info("finished getting friends")

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

func timeline(userid int64) error {
	timelineConfig, err := svc.db.GetTimelineConfig(userid)
	if err != nil {
		log.WithFields(logrus.Fields{
			"action": "timeline::GetTimelineConfig",
			"error":  err.Error(),
		}).Error("error getting timeline config")
		return err
	}

	log.WithFields(logrus.Fields{
		"action":  "timeline::Setup",
		"userid":  userid,
		"sinceid": timelineConfig.MaxID,
	}).Debug("setting up timeline")

	tweets, resp, err := svc.twitterClient.GetUserTimeline(
		&service.QueryParams{
			UserID:  userid,
			Count:   200,
			SinceID: timelineConfig.MaxID,
		},
	)
	if err != nil {
		log.WithFields(logrus.Fields{
			"action":   "timeline::GetUserTimeline",
			"response": resp.Status,
		}).Error("error getting user's timeline")
		return err
	}

	// upperID and lowerID are used to keep track of the max and min tweet IDs
	var upperID int64
	var lowerID int64

	// Loop through all the tweets.
	for t := range tweets {
		if data, err := json.Marshal(tweets[t]); err == nil {
			if err := svc.storage.Put(path.Join("timeline", strconv.FormatInt(userid, 10), tweets[t].IDStr+".json"), data); err != nil {
				logrus.WithFields(logrus.Fields{
					"action":  "timeline::Put",
					"error":   err.Error(),
					"userid":  userid,
					"tweetId": tweets[t].ID,
				}).Error("error putting timeline")
				return err
			}
		}

		// Loop through all the media entities
		for m := range tweets[t].Entities.Media {
			var url string
			if tweets[t].Entities.Media[m].MediaURLHttps != "" {
				url = tweets[t].Entities.Media[m].MediaURLHttps
			} else if tweets[t].Entities.Media[m].MediaURL != "" {
				url = tweets[t].Entities.Media[m].MediaURL
			}
			if url != "" {
				if err := svc.queue.SendEntityMessage(tweets[t].IDStr, url); err != nil {
					logrus.WithFields(logrus.Fields{
						"action":  "timeline::queue::SendMessage",
						"error":   err.Error(),
						"userid":  userid,
						"tweetId": tweets[t].ID,
					}).Error("error sending message to queue")
					fmt.Printf("Queued: %s\n", url)
				}
			}
		}

		// Calculate the max and min tweet IDs.
		if tweets[t].ID > upperID {
			//fmt.Printf("Tweet (%d) > Upper (%d), setting upperID\n", tweets[t].ID, upperID)
			upperID = tweets[t].ID
		}
		if tweets[t].ID < lowerID {
			//fmt.Printf("Tweet (%d) < Upper (%d), setting lowerID\n", tweets[t].ID, lowerID)
			lowerID = tweets[t].ID
		}
		if lowerID == 0 {
			//fmt.Printf("lowerID (%d) == 0, setting lowerID to upper ID (%d)\n", lowerID, upperID)
			lowerID = upperID
		}
	}

	if upperID > 0 {
		if err := svc.db.PutTimelineConfig(
			&database.TweetConfigQuery{
				UserID:  userid,
				SinceID: lowerID,
				MaxID:   upperID,
			},
		); err != nil {
			logrus.WithFields(logrus.Fields{
				"action":  "timeline::PutTimelineConfig",
				"error":   err.Error(),
				"userid":  userid,
				"upperID": upperID,
				"lowerID": lowerID,
			}).Error("error putting timeline config")
			return err
		}
	}

	logrus.WithFields(logrus.Fields{
		"action":  "timeline::Done!",
		"userid":  userid,
		"upperID": upperID,
		"lowerID": lowerID,
		"count":   len(tweets),
	}).Info("finished getting timeline")

	return nil
}

func user(userid int64) error {
	user, _, err := svc.twitterClient.GetUser(&service.QueryParams{UserID: userid})
	if err != nil {
		log.WithFields(logrus.Fields{
			"action": "user::GetUser",
			"userid": userid,
			"error":  err.Error(),
		}).Error("error getting user.")
		return err
	}

	log.WithFields(logrus.Fields{
		"action": "user::GetUser",
		"userid": userid,
	}).Info("got user.")

	if data, err := json.Marshal(user); err != nil {
		log.WithFields(logrus.Fields{
			"action": "user::GetUser",
			"userid": userid,
			"error":  err.Error(),
		}).Error("error marshalling user.")
		return err
	} else {
		if err := svc.storage.Put(path.Join("users", user.IDStr+".json"), data); err != nil {
			log.WithFields(logrus.Fields{
				"action": "user::GetUser",
				"userid": userid,
				"error":  err.Error(),
			}).Error("error storing user.")
			return err
		} else {
			log.WithFields(logrus.Fields{
				"action": "user::GetUser::Storage::Put",
				"userid": userid,
			}).Info("stored user.")
		}
	}
	return nil
}
