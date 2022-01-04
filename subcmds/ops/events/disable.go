package events

import (
	"github.com/sirupsen/logrus"
)

func runEventDisable() error {
	if flags.ruleName != "" {
		ruleList = []string{flags.ruleName}
	}

	for _, rule := range ruleList {

		if err := evnts.Disable(&rule); err != nil {
			log.WithFields(logrus.Fields{
				"error": err,
				"rule":  rule,
			}).Fatal("failed to disable event")
			return err
		} else {
			log.WithFields(logrus.Fields{
				"rule": rule,
			}).Info("event disabled")
		}
	}
	return nil
}
