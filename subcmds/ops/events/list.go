package events

import (
	"fmt"

	"github.com/sirupsen/logrus"
)

func runEventsList() error {
	rules, err := evnts.List()
	if err != nil {
		log.WithFields(logrus.Fields{
			"error": err,
		}).Fatal("failed to list events")
		return err
	}
	for _, rule := range rules.Rules {
		fmt.Printf("%s\n", *rule.Name)
		fmt.Printf("%s\n", *rule.Description)
		fmt.Printf("%s\n", *rule.ScheduleExpression)
		fmt.Printf("%s\n\n", *rule.State)
	}
	return nil
}
