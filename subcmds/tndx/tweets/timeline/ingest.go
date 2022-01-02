package timeline

import (
	"encoding/json"
	"fmt"

	"github.com/rmrfslashbin/tndx/pkg/queue"
	"github.com/rmrfslashbin/tndx/pkg/service"
	"github.com/sirupsen/logrus"
)

func runTimelineIngest() error {
	if flags.screenname != "" {
		if user, resp, err := svc.twitter.GetUser(&service.QueryParams{ScreenName: flags.screenname}); err != nil {
			log.WithFields(logrus.Fields{
				"action":         "runTimelineIngest::svc.twitter.GetUser",
				"err":            err,
				"responseCode":   resp.StatusCode,
				"responseStatus": resp.Status,
				"screenname":     flags.screenname,
			}).Error("error getting user's userid from screenname")
			return err
		} else {
			flags.userid = user.ID
		}
	}

	tweets, resp, err := svc.twitter.GetUserTimeline(
		&service.QueryParams{
			UserID:  flags.userid,
			Count:   flags.count,
			SinceID: flags.sinceid,
			MaxID:   flags.maxid,
		},
	)
	if err != nil {
		log.WithFields(logrus.Fields{
			"action":         "runTimelineIngest::GetUserTimeline",
			"err":            err,
			"responseCode":   resp.StatusCode,
			"responseStatus": resp.Status,
		}).Error("error getting user's timeline")
		return err
	}

	// upperID and lowerID are used to keep track of the max and min tweet IDs
	var upperID int64
	var lowerID int64

	// Loop through all the tweets.
	for t := range tweets {
		if data, err := json.Marshal(tweets[t]); err == nil {
			if opt, err := svc.kinesis.PutRecord(data); err != nil {
				log.WithFields(logrus.Fields{
					"action":  "runTimelineIngest::svc.kinesis.PutRecord",
					"error":   err,
					"tweetId": tweets[t].ID,
				}).Fatal("failed putting favorite tweet into kinesis")
			} else {
				log.WithFields(logrus.Fields{
					"action":   "runTimelineIngest::svc.kinesis.PutRecord",
					"tweetId":  tweets[t].ID,
					"recordId": *opt.RecordId,
				}).Info("put record")
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
						UserID:    flags.userid,
					},
				}); err != nil {
					logrus.WithFields(logrus.Fields{
						"action":  "runTimelineIngest::svc.queue.SendRunnerMessage",
						"error":   err.Error(),
						"userid":  flags.userid,
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

	logrus.WithFields(logrus.Fields{
		"action":  "runTimelineIngest::Done",
		"userid":  flags.userid,
		"upperID": upperID,
		"lowerID": lowerID,
		"count":   len(tweets),
	}).Info("finished getting timeline")

	return nil
}
