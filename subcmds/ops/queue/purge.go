package queue

import "github.com/sirupsen/logrus"

func runQueuePurge() error {
	if err := svc.queue.Purge(); err != nil {
		log.WithFields(logrus.Fields{
			"error": err,
		}).Fatal("failed to purge queue")
		return err
	} else {
		log.Info("queue purged")
	}
	return nil
}
