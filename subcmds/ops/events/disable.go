package events

import (
	"github.com/sirupsen/logrus"
)

func runEventDisable() error {
	if err := evnts.Disable(&flags.ruleName); err != nil {
		log.WithFields(logrus.Fields{
			"error": err,
		}).Fatal("failed to disable event")
		return err
	} else {
		log.WithFields(logrus.Fields{
			"rule": flags.ruleName,
		}).Info("event disabled")
		return nil
	}
}
