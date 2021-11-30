package subcmds

import (
	"encoding/json"
	"path"
	"strconv"

	"github.com/rmrfslashbin/tndx/pkg/database"
	"github.com/rmrfslashbin/tndx/pkg/service"
	"github.com/sirupsen/logrus"
)

func RunTimelineCmd() error {

	if flags.userid == 0 {
		user, _, err := svc.twitterClient.GetUser(&service.QueryParams{ScreenName: flags.screenname, UserID: flags.userid})
		if err != nil {
			log.WithFields(logrus.Fields{
				"action": "RunTimelineCmd::GetUser",
				"error":  err.Error(),
			}).Error("error getting user")
			return err
		}
		// Set the userid.
		flags.userid = user.ID
	}

	var maxID int64
	var sinceID int64
	if flags.maxid != 0 || flags.sinceid != 0 {
		maxID = flags.maxid
		sinceID = flags.sinceid
	} else {
		timelineConfig, err := svc.db.GetTimelineConfig(flags.userid)
		if err != nil {
			log.WithFields(logrus.Fields{
				"action": "RunTimelineCmd::GetTimelineConfig",
				"error":  err.Error(),
			}).Error("error getting timeline config")
			return err
		}

		if flags.backwards {
			maxID = timelineConfig.SinceID
		} else {
			sinceID = timelineConfig.MaxID
		}
	}

	log.WithFields(logrus.Fields{
		"action":    "RunTimelineCmd::Setup",
		"userid":    flags.userid,
		"maxid":     maxID,
		"sinceid":   sinceID,
		"backwards": flags.backwards,
	}).Debug("setting up timeline")

	tweets, resp, err := svc.twitterClient.GetUserTimeline(
		&service.QueryParams{
			UserID:  flags.userid,
			MaxID:   maxID,
			Count:   200,
			SinceID: sinceID,
		},
	)
	if err != nil {
		log.WithFields(logrus.Fields{
			"action":   "RunTimelineCmd::GetUserTimeline",
			"response": resp.Status,
		}).Error("error getting user's timeline")
		return err
	}

	// upperID and lowerID are used to keep track of the max and min tweet IDs
	var upperID int64
	var lowerID int64

	// Loop through all the tweets.
	for t := range tweets {
		if data, err := json.MarshalIndent(tweets[t], "", "  "); err == nil {
			if err := svc.storage.Put(path.Join(strconv.FormatInt(flags.userid, 10), "timeline", tweets[t].IDStr), data); err != nil {
				logrus.WithFields(logrus.Fields{
					"action":  "RunTimelineCmd::Put",
					"error":   err.Error(),
					"userid":  flags.userid,
					"tweetId": tweets[t].ID,
				}).Error("error putting timeline")
				return err
			}
		}

		// Loop through all the media entities
		/*
			for m := range tweets[t].Entities.Media {
				var url string
				if tweets[t].Entities.Media[m].MediaURLHttps != "" {
					url = tweets[t].Entities.Media[m].MediaURLHttps
				} else if tweets[t].Entities.Media[m].MediaURL != "" {
					url = tweets[t].Entities.Media[m].MediaURL
				}
				if url != "" {
					if err := svc.queue.SendMessage(flags.userid, url); err != nil {
						logrus.WithFields(logrus.Fields{
							"action":  "RunTimelineCmd::queue::SendMessage",
							"error":   err.Error(),
							"userid":  flags.userid,
							"tweetId": tweets[t].ID,
						}).Error("error sending message to queue")
						fmt.Printf("Queued: %s\n", url)
					}
				}
			}
		*/

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
				UserID:  flags.userid,
				SinceID: lowerID,
				MaxID:   upperID,
			},
		); err != nil {
			logrus.WithFields(logrus.Fields{
				"action":  "RunTimelineCmd::PutTimelineConfig",
				"error":   err.Error(),
				"userid":  flags.userid,
				"upperID": upperID,
				"lowerID": lowerID,
			}).Error("error putting timeline config")
			return err
		}
	}

	logrus.WithFields(logrus.Fields{
		"action":  "RunTimelineCmd::Done!",
		"userid":  flags.userid,
		"upperID": upperID,
		"lowerID": lowerID,
		"count":   len(tweets),
	}).Info("finished getting timeline")

	return nil
}
