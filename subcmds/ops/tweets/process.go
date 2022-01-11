package tweets

import (
	"github.com/rmrfslashbin/tndx/pkg/queue"
	"github.com/sirupsen/logrus"
)

func RunTweetsProcess() error {
	bootstrap.Function = "get_tweet"
	for _, tweetid := range flags.tweetids {
		if err := svc.queue.SendRunnerMessage(&queue.SendMessage{
			Bootstrap: bootstrap,
			Message: &queue.ProcessorMessage{
				TweetID: tweetid,
			},
		}); err != nil {
			logrus.WithFields(logrus.Fields{
				"action":  "timeline::queue::SendRunnerMessage",
				"error":   err.Error(),
				"tweetId": tweetid,
			}).Error("error sending message to queue")
		} else {
			logrus.WithFields(logrus.Fields{
				"action":   "RunTweetsProcess::queue::SendRunnerMessage",
				"function": bootstrap.Function,
				"tweetId":  tweetid,
			}).Info("message sent to queue")
		}
	}

	return nil
}
