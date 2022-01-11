package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/rmrfslashbin/tndx/pkg/database"
	"github.com/rmrfslashbin/tndx/pkg/kinesis"
	"github.com/rmrfslashbin/tndx/pkg/queue"
	"github.com/rmrfslashbin/tndx/pkg/service"
	"github.com/rmrfslashbin/tndx/pkg/ssmparams"
	"github.com/rmrfslashbin/tndx/pkg/storage"
	"github.com/sirupsen/logrus"
)

// service stores drivers and clients
type services struct {
	twitterClient *service.Config
	storage       *storage.S3Storage
	db            *database.DDBDriver
	queue         *queue.Config
	kinesis       *kinesis.Config
}

var (
	aws_region string
	log        *logrus.Logger
	svc        *services
)

func init() {
	log = logrus.New()
	log.SetLevel(logrus.InfoLevel)
	log.SetFormatter(&logrus.JSONFormatter{})
	aws_region = os.Getenv("AWS_REGION")
	svc = &services{}
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
	log.WithFields(logrus.Fields{
		"action": "handler",
		"event":  sqsEvent,
	}).Info("starting handler")
	for _, message := range sqsEvent.Records {
		messageBody := &queue.ProcessorMessage{}
		if err := json.Unmarshal([]byte(message.Body), messageBody); err != nil {
			logrus.WithFields(logrus.Fields{
				"action": "handler::Unmarshal",
				"error":  err.Error(),
				"body":   message.Body,
			}).Error("error unmarshalling message")
			return err
		}

		if _, ok := message.MessageAttributes["function"]; !ok {
			return errors.New("function is required")
		}
		if _, ok := message.MessageAttributes["ddb_table_prefix"]; !ok {
			return errors.New("ddb table prefix is required")
		}
		if _, ok := message.MessageAttributes["delivery_stream"]; !ok {
			return errors.New("delivery stream is required")
		}
		if _, ok := message.MessageAttributes["sqs_runner_url"]; !ok {
			return errors.New("sqs runner url is required")
		}
		if _, ok := message.MessageAttributes["s3_bucket"]; !ok {
			return errors.New("s3 bucket is required")
		}
		if _, ok := message.MessageAttributes["twitter_api_key"]; !ok {
			return errors.New("twitter api key is required")
		}
		if _, ok := message.MessageAttributes["twitter_api_secret"]; !ok {
			return errors.New("twitter api secret is required")
		}

		bootstrap := &queue.Bootstrap{
			Function:         *message.MessageAttributes["function"].StringValue,
			DDBTablePrefix:   *message.MessageAttributes["ddb_table_prefix"].StringValue,
			DeliveryStream:   *message.MessageAttributes["delivery_stream"].StringValue,
			SQSRunnerURL:     *message.MessageAttributes["sqs_runner_url"].StringValue,
			S3Bucket:         *message.MessageAttributes["s3_bucket"].StringValue,
			TwitterAPIKey:    *message.MessageAttributes["twitter_api_key"].StringValue,
			TwitterAPISecret: *message.MessageAttributes["twitter_api_secret"].StringValue,
		}

		if bootstrap.Function == "" {
			return errors.New("function is required")
		}
		if bootstrap.DDBTablePrefix == "" {
			return errors.New("ddb table prefix is required")
		}
		if bootstrap.DeliveryStream == "" {
			return errors.New("delivery stream is required")
		}
		if bootstrap.SQSRunnerURL == "" {
			return errors.New("sqs runner url is required")
		}
		if bootstrap.S3Bucket == "" {
			return errors.New("s3 bucket is required")
		}
		if bootstrap.TwitterAPIKey == "" {
			return errors.New("twitter api key is required")
		}
		if bootstrap.TwitterAPISecret == "" {
			return errors.New("twitter api secret is required")
		}

		params := ssmparams.NewSSMParams(
			ssmparams.SetRegion(aws_region),
			ssmparams.SetLogger(log),
		)

		outputs, err := params.GetParams([]string{
			bootstrap.DDBTablePrefix,
			bootstrap.DeliveryStream,
			bootstrap.SQSRunnerURL,
			bootstrap.S3Bucket,
			bootstrap.TwitterAPIKey,
			bootstrap.TwitterAPISecret,
		})

		if err != nil {
			log.WithFields(logrus.Fields{
				"action":    "getParams",
				"error":     err.Error(),
				"bootstrap": bootstrap,
			}).Error("error getting parameters.")
			return err
		}

		if len(outputs.InvalidParameters) > 0 {
			log.WithFields(logrus.Fields{
				"invalid_parameters": outputs.InvalidParameters,
			}).Error("invalid parameters")
			return errors.New("invalid parameters")
		}

		svc.db = database.NewDDB(
			database.SetDDBLogger(log),
			database.SetDDBRegion(aws_region),
			database.SetDDBTablePrefix(outputs.Params[bootstrap.DDBTablePrefix].(string)),
		)

		svc.queue = queue.NewSQS(
			queue.SetLogger(log),
			queue.SetSQSURL(outputs.Params[bootstrap.SQSRunnerURL].(string)),
		)

		svc.storage = storage.NewS3Storage(
			storage.SetS3Bucket(outputs.Params[bootstrap.S3Bucket].(string)),
			storage.SetS3Region(aws_region),
			storage.SetLogger(log),
		)

		svc.twitterClient = service.New(
			service.SetConsumerKey(outputs.Params[bootstrap.TwitterAPIKey].(string)),
			service.SetConsumerSecret(outputs.Params[bootstrap.TwitterAPISecret].(string)),
			service.SetLogger(log),
		)

		svc.kinesis = kinesis.NewFirehose(
			kinesis.SetRegion(aws_region),
			kinesis.SetLogger(log),
			kinesis.SetDeliveryStream(outputs.Params[bootstrap.DeliveryStream].(string)),
		)

		switch bootstrap.Function {
		case "entities":
			if err := entities(&messageBody.UserID, &messageBody.TweetID, &messageBody.EntityURL); err != nil {
				logrus.WithFields(logrus.Fields{
					"function": "entities",
					"error":    err,
				}).Error("function failed")
				return err
			}

		case "favorites":
			if err := favorites(messageBody.UserID, bootstrap); err != nil {
				logrus.WithFields(logrus.Fields{
					"function": "favorites",
					"error":    err,
				}).Error("function failed")
				return err
			}

		case "followers":
			if err := followers(messageBody.UserID); err != nil {
				logrus.WithFields(logrus.Fields{
					"function": "followers",
					"error":    err,
				}).Error("function failed")
				return err
			}

		case "friends":
			if err := friends(messageBody.UserID); err != nil {
				logrus.WithFields(logrus.Fields{
					"function": "friends",
					"error":    err,
				}).Error("function failed")
				return err
			}

		case "get_tweet":
			tweetId, err := strconv.ParseInt(messageBody.TweetID, 10, 64)
			if err != nil {
				logrus.WithFields(logrus.Fields{
					"function": "get_tweet",
					"error":    err,
				}).Error("unable to parse tweet id to int64")
				return err
			}
			if err := getTweet(tweetId, bootstrap); err != nil {
				logrus.WithFields(logrus.Fields{
					"function": "get_tweet",
					"error":    err,
				}).Error("function failed")
				return err
			}

		case "timeline":
			if err := timeline(messageBody.UserID, bootstrap); err != nil {
				logrus.WithFields(logrus.Fields{
					"function": "timeline",
					"error":    err,
				}).Error("function failed")
				return err
			}

		case "user":
			if err := user(messageBody.UserID); err != nil {
				logrus.WithFields(logrus.Fields{
					"function": "user",
					"error":    err,
				}).Error("function failed")
				return err
			}

		default:
			logrus.WithFields(logrus.Fields{
				"function": bootstrap.Function,
			}).Error("invalid function; should be one of user, friend, followers, favorites, timeline, entities")
			return errors.New("invalid function; should be one of user, friend, followers, favorites, timeline, entities")
		}
	}

	return nil
}

func entities(userId *int64, tweetId *string, entityURL *string) error {
	resp, err := http.Get(*entityURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	filenameParts := strings.Split(*entityURL, "/")
	key := fmt.Sprintf("media/%d/%s/%s", *userId, *tweetId, filenameParts[len(filenameParts)-1])

	if err := svc.storage.PutStream(key, resp.Body); err != nil {
		return err
	}

	log.WithFields(logrus.Fields{
		"action":    "entites",
		"userid":    userId,
		"tweetId":   tweetId,
		"entityURL": entityURL,
	}).Info("fetched and put entity")
	return nil
}

func favorites(userid int64, bootstrap *queue.Bootstrap) error {
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
		if resp.StatusCode == 429 {
			log.WithFields(logrus.Fields{
				"action":         "favorites::GetUserFavorites",
				"error":          err,
				"responsestatus": resp.Header,
			}).Error("rate limit exceeded getting user's favorites")
			return err
		} else {
			log.WithFields(logrus.Fields{
				"action":         "favorites",
				"responseCode":   resp.StatusCode,
				"responseStatus": resp.Status,
			}).Error("error getting user's favorites")
			return err
		}
	}

	// upperID and lowerID are used to keep track of the max and min tweet IDs
	var upperID int64
	var lowerID int64

	listOfTweets := make([]*database.UserToTweetLink, len(tweets))

	// Loop through all the tweets.
	for t := range tweets {
		listOfTweets[t] = &database.UserToTweetLink{UserID: userid, TweetID: tweets[t].ID}
		if data, err := json.Marshal(tweets[t]); err == nil {
			if opt, err := svc.kinesis.PutRecord(data); err != nil {
				log.WithFields(logrus.Fields{
					"error":   err,
					"tweetId": tweets[t].ID,
				}).Fatal("failed putting favorite tweet into kinesis")
			} else {
				log.WithFields(logrus.Fields{
					"tweetId":  tweets[t].ID,
					"recordId": *opt.RecordId,
				}).Info("put record")
			}
		}

		// check for RetweetedStatus
		if tweets[t].RetweetedStatus != nil {
			bootstrap.Function = "get_tweet"
			if err := svc.queue.SendRunnerMessage(&queue.SendMessage{
				Bootstrap: bootstrap,
				Message: &queue.ProcessorMessage{
					TweetID: tweets[t].RetweetedStatus.IDStr,
				},
			}); err != nil {
				logrus.WithFields(logrus.Fields{
					"action":  "timeline::queue::SendRunnerMessage",
					"error":   err.Error(),
					"tweetId": tweets[t].ID,
				}).Error("error sending message to queue")
			}
		}

		// check for quoted_status_id
		if tweets[t].QuotedStatusIDStr != "" {
			bootstrap.Function = "get_tweet"
			if err := svc.queue.SendRunnerMessage(&queue.SendMessage{
				Bootstrap: bootstrap,
				Message: &queue.ProcessorMessage{
					TweetID: tweets[t].QuotedStatusIDStr,
				},
			}); err != nil {
				logrus.WithFields(logrus.Fields{
					"action":  "favorites::queue::SendRunnerMessage",
					"error":   err.Error(),
					"tweetId": tweets[t].ID,
				}).Error("error sending message to queue")
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
				bootstrap.Function = "entities"
				if err := svc.queue.SendRunnerMessage(&queue.SendMessage{
					Bootstrap: bootstrap,
					Message: &queue.ProcessorMessage{
						TweetID:   tweets[t].IDStr,
						EntityURL: url,
						UserID:    userid,
					},
				}); err != nil {
					logrus.WithFields(logrus.Fields{
						"action":  "favorites::queue::SendRunnerMessage",
						"error":   err.Error(),
						"userid":  userid,
						"tweetId": tweets[t].ID,
					}).Error("error sending message to queue")
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

	if err := svc.db.PutFavorites(listOfTweets); err != nil {
		logrus.WithFields(logrus.Fields{
			"action": "favorites::PutFavorites",
			"error":  err.Error(),
		}).Error("error putting favorites")
		return err
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
		if resp.StatusCode == 429 {
			log.WithFields(logrus.Fields{
				"action":         "followers::GetUserFollowers",
				"error":          err,
				"responsestatus": resp.Header,
			}).Error("rate limit exceeded getting user's followers")
			return err
		} else {
			log.WithFields(logrus.Fields{
				"action": "followers::GetUserFollowers",
				"error":  err,
			}).Error("error getting user's followers")
			return err
		}
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

	listOfFollowers := make([]*database.UserToFollowerLink, len(followers.Users))

	// Save the users.
	for f := range followers.Users {
		listOfFollowers[f] = &database.UserToFollowerLink{UserID: userid, FollowerID: followers.Users[f].ID}
		if data, err := json.Marshal(followers.Users[f]); err == nil {
			if err := svc.storage.Put(path.Join("users", followers.Users[f].IDStr+".json"), data); err != nil {
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
	if err := svc.db.PutFollowers(listOfFollowers); err != nil {
		logrus.WithFields(logrus.Fields{
			"action": "followers::PutFollowers",
			"error":  err.Error(),
		}).Error("error putting followers")
		return err
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
		if resp.StatusCode == 429 {
			log.WithFields(logrus.Fields{
				"action":         "friends::GetUserFriends",
				"error":          err,
				"responsestatus": resp.Header,
			}).Error("rate limit exceeded getting user's friends")
			return err
		} else {
			log.WithFields(logrus.Fields{
				"action":   "friends::GetUserFriends",
				"response": resp.Status,
			}).Error("error getting user's friends")
			return err
		}
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

	listOfFriends := make([]*database.UserToFriendLink, len(friends.Users))

	// Save the users.
	for f := range friends.Users {
		listOfFriends[f] = &database.UserToFriendLink{UserID: userid, FriendID: friends.Users[f].ID}
		if data, err := json.Marshal(friends.Users[f]); err == nil {
			if err := svc.storage.Put(path.Join("users", friends.Users[f].IDStr+".json"), data); err != nil {
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
	if err := svc.db.PutFriends(listOfFriends); err != nil {
		logrus.WithFields(logrus.Fields{
			"action": "friends::PutFriends",
			"error":  err.Error(),
		}).Error("error putting friends")
		return err
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

func getTweet(tweetId int64, bootstrap *queue.Bootstrap) error {
	tweets, resp, err := svc.twitterClient.LookupTweets([]int64{tweetId})
	if err != nil {
		if resp.StatusCode == 429 {
			log.WithFields(logrus.Fields{
				"action":         "getTweet::svc.twitterClient.LookupTweets",
				"error":          err,
				"responsestatus": resp.Header,
				"tweetId":        tweetId,
			}).Error("rate limit exceeded getting tweet")
			return err
		} else {
			log.WithFields(logrus.Fields{
				"action":         "getTweet::svc.twitterClient.LookupTweets",
				"responseCode":   resp.StatusCode,
				"responseStatus": resp.Status,
				"tweetId":        tweetId,
				"error":          err.Error(),
			}).Error("error getting tweet")
			return err
		}
	}

	// Loop through all the tweets.
	for t := range tweets {
		log.WithFields(logrus.Fields{
			"action":  "getTweet::LookupTweets",
			"tweetId": tweets[t],
		}).Info("base tweet")
		if data, err := json.Marshal(tweets[t]); err == nil {
			if opt, err := svc.kinesis.PutRecord(data); err != nil {
				log.WithFields(logrus.Fields{
					"action":  "getTweet::svc.kinesis.PutRecord",
					"error":   err,
					"tweetId": tweets[t].ID,
				}).Fatal("failed putting tweet into kinesis")
			} else {
				log.WithFields(logrus.Fields{
					"tweetId":  tweets[t].ID,
					"recordId": *opt.RecordId,
				}).Info("put record")
			}
		}

		// check for RetweetedStatus
		if tweets[t].RetweetedStatus != nil {
			bootstrap.Function = "get_tweet"
			if err := svc.queue.SendRunnerMessage(&queue.SendMessage{
				Bootstrap: bootstrap,
				Message: &queue.ProcessorMessage{
					TweetID: tweets[t].RetweetedStatus.IDStr,
				},
			}); err != nil {
				logrus.WithFields(logrus.Fields{
					"action":  "getTweet::svc.queue.SendRunnerMessage::get_tweet::RetweetedStatus",
					"error":   err.Error(),
					"tweetId": tweets[t].ID,
				}).Error("error sending message to queue")
			}
		}

		// check for quoted_status_id
		if tweets[t].QuotedStatusIDStr != "" {
			bootstrap.Function = "get_tweet"
			if err := svc.queue.SendRunnerMessage(&queue.SendMessage{
				Bootstrap: bootstrap,
				Message: &queue.ProcessorMessage{
					TweetID: tweets[t].QuotedStatusIDStr,
				},
			}); err != nil {
				logrus.WithFields(logrus.Fields{
					"action":  "getTweet::svc.queue.SendRunnerMessage::get_tweet::QuotedStatusIDStr",
					"error":   err.Error(),
					"tweetId": tweets[t].ID,
				}).Error("error sending message to queue")
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
				bootstrap.Function = "entities"
				if err := svc.queue.SendRunnerMessage(&queue.SendMessage{
					Bootstrap: bootstrap,
					Message: &queue.ProcessorMessage{
						TweetID:   tweets[t].IDStr,
						EntityURL: url,
						UserID:    tweets[t].User.ID,
					},
				}); err != nil {
					logrus.WithFields(logrus.Fields{
						"action":  "getTweet::svc.queue.SendRunnerMessage::entities",
						"error":   err.Error(),
						"userid":  tweets[t].User.ID,
						"tweetId": tweets[t].ID,
					}).Error("error sending message to queue")
				}
			}
		}
	}

	logrus.WithFields(logrus.Fields{
		"action":  "get_tweet::Done!",
		"tweetId": tweetId,
		"count":   len(tweets),
	}).Info("finished getting tweet")

	return nil
}

func timeline(userid int64, bootstrap *queue.Bootstrap) error {
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
		if resp.StatusCode == 429 {
			log.WithFields(logrus.Fields{
				"action":         "followers::GetUserFollowers",
				"error":          err,
				"responsestatus": resp.Header,
			}).Error("rate limit exceeded getting user's followers")
			log.WithFields(logrus.Fields{
				"action": "followers::GetUserFollowers",
				"error":  err,
			}).Error("rate limit exceeded getting user's followers")
			return nil
		} else {
			log.WithFields(logrus.Fields{
				"action":         "timeline::GetUserTimeline",
				"responseCode":   resp.StatusCode,
				"responseStatus": resp.Status,
			}).Error("error getting user's timeline")
			return err
		}
	}

	// upperID and lowerID are used to keep track of the max and min tweet IDs
	var upperID int64
	var lowerID int64

	// Loop through all the tweets.
	for t := range tweets {
		log.WithFields(logrus.Fields{
			"action": "timeline::GetUserTimeline",
			"tweet":  tweets[t],
		}).Info("base tweet")
		if data, err := json.Marshal(tweets[t]); err == nil {
			if opt, err := svc.kinesis.PutRecord(data); err != nil {
				log.WithFields(logrus.Fields{
					"error":   err,
					"tweetId": tweets[t].ID,
				}).Fatal("failed putting favorite tweet into kinesis")
			} else {
				log.WithFields(logrus.Fields{
					"tweetId":  tweets[t].ID,
					"recordId": *opt.RecordId,
				}).Info("put record")
			}
		}

		// check for RetweetedStatus
		if tweets[t].RetweetedStatus != nil {
			bootstrap.Function = "get_tweet"
			if err := svc.queue.SendRunnerMessage(&queue.SendMessage{
				Bootstrap: bootstrap,
				Message: &queue.ProcessorMessage{
					TweetID: tweets[t].RetweetedStatus.IDStr,
				},
			}); err != nil {
				logrus.WithFields(logrus.Fields{
					"action":  "timeline::queue::SendRunnerMessage",
					"error":   err.Error(),
					"tweetId": tweets[t].ID,
				}).Error("error sending message to queue")
			}
		}

		// check for quoted_status_id
		if tweets[t].QuotedStatusIDStr != "" {
			bootstrap.Function = "get_tweet"
			if err := svc.queue.SendRunnerMessage(&queue.SendMessage{
				Bootstrap: bootstrap,
				Message: &queue.ProcessorMessage{
					TweetID: tweets[t].QuotedStatusIDStr,
				},
			}); err != nil {
				logrus.WithFields(logrus.Fields{
					"action":  "timeline::queue::SendRunnerMessage",
					"error":   err.Error(),
					"tweetId": tweets[t].ID,
				}).Error("error sending message to queue")
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
				bootstrap.Function = "entities"
				if err := svc.queue.SendRunnerMessage(&queue.SendMessage{
					Bootstrap: bootstrap,
					Message: &queue.ProcessorMessage{
						TweetID:   tweets[t].IDStr,
						EntityURL: url,
						UserID:    userid,
					},
				}); err != nil {
					logrus.WithFields(logrus.Fields{
						"action":  "timeline::queue::SendMessage",
						"error":   err.Error(),
						"userid":  userid,
						"tweetId": tweets[t].ID,
					}).Error("error sending message to queue")
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
	user, resp, err := svc.twitterClient.GetUser(&service.QueryParams{UserID: userid})
	if err != nil {
		if resp.StatusCode == 429 {
			log.WithFields(logrus.Fields{
				"action":         "followers::GetUserFollowers",
				"error":          err,
				"responsestatus": resp.Header,
			}).Error("rate limit exceeded getting user's followers")
			log.WithFields(logrus.Fields{
				"action": "followers::GetUserFollowers",
				"error":  err,
			}).Error("rate limit exceeded getting user's followers")
			return nil
		} else {
			log.WithFields(logrus.Fields{
				"action": "user::GetUser",
				"userid": userid,
				"error":  err.Error(),
			}).Error("error getting user.")
			return err
		}
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
