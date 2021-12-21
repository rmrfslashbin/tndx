package runner

import (
	"fmt"

	"github.com/rmrfslashbin/tndx/pkg/database"
	"github.com/rmrfslashbin/tndx/pkg/service"
	"github.com/sirupsen/logrus"
)

func RunRunnerList() error {
	if flags.userid == 0 && flags.screenname != "" && !flags.all {
		user, _, err := svc.twitterClient.GetUser(&service.QueryParams{ScreenName: flags.screenname, UserID: flags.userid})
		if err != nil {
			log.WithFields(logrus.Fields{
				"action": "RunRunnerList::GetUser",
				"error":  err.Error(),
			}).Error("error getting user")
			return err
		}
		// Set the userid.
		flags.userid = user.ID
	}

	if flags.userid == 0 && !flags.all {
		log.WithFields(logrus.Fields{
			"action": "RunRunnerList::GetUser",
			"error":  "no userid or screenname provided/could not be resolved",
		}).Fatal("error getting user")
	}

	res, err := svc.db.GetRunnerUsers(&database.RunnerItem{
		RunnerName: flags.runner,
		UserID:     flags.userid,
	})
	if err != nil {
		log.WithFields(logrus.Fields{
			"action": "RunRunnerList",
			"error":  err.Error(),
		}).Error("error getting runner users")
		return err
	}

	if len(res) == 0 {
		log.WithFields(logrus.Fields{
			"runner":     flags.runner,
			"userid":     flags.userid,
			"screenname": flags.screenname,
		}).Info("no users found for runner")
	}

	for _, user := range res {
		logrus.WithFields(logrus.Fields{
			"action":      "RunRunnerList",
			"userid":      user.UserID,
			"flags":       user.Flags,
			"newFlagsBin": fmt.Sprintf("%08b", user.Flags),
			"favorites":   database.Has(user.Flags, database.F_favorites),
			"followers":   database.Has(user.Flags, database.F_followers),
			"friends":     database.Has(user.Flags, database.F_friends),
			"timeline":    database.Has(user.Flags, database.F_timeline),
			"user":        database.Has(user.Flags, database.F_user),
		}).Info("flags")
	}
	return nil
}
