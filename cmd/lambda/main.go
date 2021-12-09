package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"path"
	"strconv"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/rmrfslashbin/tndx/pkg/database"
	"github.com/rmrfslashbin/tndx/pkg/queue"
	"github.com/rmrfslashbin/tndx/pkg/service"
	"github.com/rmrfslashbin/tndx/pkg/storage"
	"github.com/sirupsen/logrus"
)

type SSMParams struct {
	EntityQueue      string `json:"entity_queue"`
	S3Bucket         string `json:"s3_bucket"`
	S3Region         string `json:"s3_region"`
	DDBTable         string `json:"ddb_table"`
	DDBRegion        string `json:"ddb_region"`
	TwitterAPIKey    string `json:"twitter_api_key"`
	TwitterAPISecret string `json:"twitter_api_secret"`
}

type Bootstrap struct {
	Function  string    `json:"function"` // user, friends, followers, favorties, timeline, entities
	Loglevel  string    `json:"loglevel"` // error, warn, info, debug, trace
	UserID    int64     `json:"userid"`
	SSMParams SSMParams `json:"ssm_params"`
}

type Response struct {
	Message string `json:"message"`
}

// service stores drivers and clients
type services struct {
	twitterClient *service.Config
	storage       storage.StorageDriver
	db            database.DatabaseDriver
	queue         *queue.Config
}

var log *logrus.Logger
var svc services

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
	lambda.Start(HandleLambdaEvent)
}

func HandleLambdaEvent(event Bootstrap) (Response, error) {
	message := ""

	// Set log level
	switch event.Loglevel {
	case "debug":
		log.SetLevel(logrus.DebugLevel)

	case "info":
		log.SetLevel(logrus.InfoLevel)

	case "warn":
		log.SetLevel(logrus.WarnLevel)

	case "error":
		log.SetLevel(logrus.ErrorLevel)

	case "trace":
		log.SetLevel(logrus.TraceLevel)

	default:
		log.SetLevel(logrus.InfoLevel)
	}

	if event.SSMParams.EntityQueue == "" {
		return Response{Message: "SSMParams.EntityQueue is required"}, errors.New("SSMParams.EntityQueue is required")
	}
	if event.SSMParams.S3Bucket == "" {
		return Response{Message: "SSMParams.S3Bucket is required"}, errors.New("SSMParams.S3Bucket is required")
	}
	if event.SSMParams.S3Region == "" {
		return Response{Message: "SSMParams.S3Region is required"}, errors.New("SSMParams.S3Region is required")
	}
	if event.SSMParams.DDBTable == "" {
		return Response{Message: "SSMParams.DDBTable is required"}, errors.New("SSMParams.DDBTable is required")
	}
	if event.SSMParams.DDBRegion == "" {
		return Response{Message: "SSMParams.DDBRegion is required"}, errors.New("SSMParams.DDBRegion is required")
	}
	if event.SSMParams.TwitterAPIKey == "" {
		return Response{Message: "SSMParams.TwitterAPIKey is required"}, errors.New("SSMParams.TwitterAPIKey is required")
	}
	if event.SSMParams.TwitterAPISecret == "" {
		return Response{Message: "SSMParams.TwitterAPISecret is required"}, errors.New("SSMParams.TwitterAPISecret is required")
	}

	outputs, err := getParams([]*string{
		&event.SSMParams.EntityQueue,
		&event.SSMParams.S3Bucket,
		&event.SSMParams.S3Region,
		&event.SSMParams.DDBTable,
		&event.SSMParams.DDBRegion,
		&event.SSMParams.TwitterAPIKey,
		&event.SSMParams.TwitterAPISecret,
	})

	if err != nil {
		logrus.WithFields(logrus.Fields{
			"function": event.Function,
			"action":   "ssmparams::GetParams",
			"error":    err.Error(),
		}).Error("ssmparams.GetParams failed to get params")
		return Response{}, err
	}

	svc.queue = queue.NewSQS(
		queue.SetLogger(log),
		queue.SetSQSURL(outputs[event.SSMParams.EntityQueue].(string)),
		queue.SetS3Bucket(outputs[event.SSMParams.S3Bucket].(string)),
	)

	svc.db = database.NewDDB(
		database.SetDDBLogger(log),
		database.SetDDBTable(outputs[event.SSMParams.DDBTable].(string)),
		database.SetDDBRegion(outputs[event.SSMParams.DDBRegion].(string)),
	)

	svc.storage = storage.NewS3Storage(
		storage.SetS3Bucket(outputs[event.SSMParams.S3Bucket].(string)),
		storage.SetS3Region(outputs[event.SSMParams.S3Region].(string)),
	)

	svc.twitterClient = service.New(
		service.SetConsumerKey(outputs[event.SSMParams.TwitterAPIKey].(string)),
		service.SetConsumerSecret(outputs[event.SSMParams.TwitterAPISecret].(string)),
		service.SetLogger(log),
	)

	switch event.Function {
	case "favorites":
		if err := favorites(event.UserID); err != nil {
			logrus.WithFields(logrus.Fields{
				"function": event.Function,
				"error":    err,
			}).Error("function failed")
			return Response{}, err
		}
		message = fmt.Sprintf("finished fetching favorites for userid: %d", event.UserID)

	case "followers":
		if err := followers(event.UserID); err != nil {
			logrus.WithFields(logrus.Fields{
				"function": event.Function,
				"error":    err,
			}).Error("function failed")
			return Response{}, err
		}
		message = fmt.Sprintf("finished fetching followers for userid: %d", event.UserID)

	case "friends":
		if err := friends(event.UserID); err != nil {
			logrus.WithFields(logrus.Fields{
				"function": event.Function,
				"error":    err,
			}).Error("function failed")
			return Response{}, err
		}
		message = fmt.Sprintf("finished fetching friends for userid: %d", event.UserID)

	case "timeline":
		if err := timeline(event.UserID); err != nil {
			logrus.WithFields(logrus.Fields{
				"function": event.Function,
				"error":    err,
			}).Error("function failed")
			return Response{}, err
		}
		message = fmt.Sprintf("finished fetching timeline for userid: %d", event.UserID)

	case "user":
		if err := user(event.UserID); err != nil {
			logrus.WithFields(logrus.Fields{
				"function": event.Function,
				"error":    err,
			}).Error("function failed")
			return Response{}, err
		}
		message = fmt.Sprintf("finished fetching user for userid: %d", event.UserID)

	case "entities":
		message = "entities"

	default:
		logrus.WithFields(logrus.Fields{
			"function": event.Function,
		}).Error("invalid function; should be one of user, friend, followers, favorites, timeline, entities")
		return Response{}, errors.New("invalid function; should be one of user, friend, followers, favorites, timeline, entities")
	}

	logrus.WithFields(logrus.Fields{
		"message":  message,
		"function": event.Function,
	}).Info("Lambda function triggered")

	return Response{Message: message}, nil
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
