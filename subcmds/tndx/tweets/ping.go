package tweets

import (
	"github.com/sirupsen/logrus"
)

func runTweetsPing() error {
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

	foundTweets := make(map[int64]bool, len(tweets))
	for _, tweet := range tweets {
		foundTweets[tweet.ID] = true
	}

	for _, tweetids := range tweetIDs {
		if _, ok := foundTweets[tweetids]; ok {
			log.WithFields(logrus.Fields{
				"tweetid": tweetids,
			}).Info("found tweet")
		} else {
			log.WithFields(logrus.Fields{
				"tweetids": tweetids,
			}).Error("tweet not found")
		}
	}
	return nil
}
