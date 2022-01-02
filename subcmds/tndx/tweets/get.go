package tweets

import (
	"encoding/json"
	"os"

	"github.com/davecgh/go-spew/spew"
	"github.com/rmrfslashbin/tndx/pkg/service"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

func runTweetsGet() error {
	tweetIDs := DedupInt64Slice(flags.tweetids)
	tweets, resp, err := svc.twitter.LookupTweets(tweetIDs)
	if err != nil {
		log.WithFields(logrus.Fields{
			"error":  err,
			"code":   resp.StatusCode,
			"status": resp.Status,
			"tweets": tweetIDs,
		}).Error("error looking up tweets")
		return err
	}
	log.WithFields(logrus.Fields{
		"count": len(tweets),
	}).Info("tweets returned")

	for t := range tweets {
		tweets[t].CreatedAt, _ = service.FixTwitterTime(tweets[t].CreatedAt)
		tweets[t].User.CreatedAt, _ = service.FixTwitterTime(tweets[t].User.CreatedAt)
		if tweets[t].RetweetedStatus != nil {
			tweets[t].RetweetedStatus.CreatedAt, _ = service.FixTwitterTime(tweets[t].RetweetedStatus.CreatedAt)
			tweets[t].RetweetedStatus.User.CreatedAt, _ = service.FixTwitterTime(tweets[t].RetweetedStatus.User.CreatedAt)
		}
	}

	if flags.json {
		if data, err := json.Marshal(tweets); err != nil {
			log.WithFields(logrus.Fields{
				"error": err,
			}).Error("error marshalling tweets to json")
			return err
		} else {
			os.Stdout.Write(data)
		}
	} else if flags.yaml {
		if data, err := yaml.Marshal(tweets); err != nil {
			log.WithFields(logrus.Fields{
				"error": err,
			}).Error("error marshalling tweets to yaml")
			return err
		} else {
			os.Stdout.Write(data)
		}
	} else {
		spew.Dump(tweets)
	}
	return nil
}
