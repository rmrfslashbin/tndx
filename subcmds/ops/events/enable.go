package events

import (
	"github.com/sirupsen/logrus"
)

func runEventEnable() error {
	if err := evnts.Enable(&flags.ruleName); err != nil {
		log.WithFields(logrus.Fields{
			"error": err,
		}).Fatal("failed to enable event")
		return err
	} else {
		log.WithFields(logrus.Fields{
			"rule": flags.ruleName,
		}).Info("event enabled")
		return nil
	}
}
