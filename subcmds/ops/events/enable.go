package events

import (
	"github.com/sirupsen/logrus"
)

func runEventEnable() error {
	if flags.ruleName != "" {
		ruleList = []string{flags.ruleName}
	}

	for _, rule := range ruleList {
		if err := evnts.Enable(&rule); err != nil {
			log.WithFields(logrus.Fields{
				"error": err,
				"rule":  rule,
			}).Fatal("failed to enable event")
			return err
		} else {
			log.WithFields(logrus.Fields{
				"rule": rule,
			}).Info("event enabled")
		}
	}
	return nil
}
