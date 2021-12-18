package subcmds

import (
	"fmt"

	"github.com/rmrfslashbin/tndx/pkg/database"
	"github.com/rmrfslashbin/tndx/pkg/service"
	"github.com/sirupsen/logrus"
)

func RunDDBRunnerList() error {
	if flags.userid == 0 && flags.screenname != "" {
		user, _, err := svc.twitterClient.GetUser(&service.QueryParams{ScreenName: flags.screenname, UserID: flags.userid})
		if err != nil {
			log.WithFields(logrus.Fields{
				"action": "RunTimelineCmd::GetUser",
				"error":  err.Error(),
			}).Error("error getting user")
			return err
		}
		// Set the userid.
		flags.userid = user.ID
	}

	res, err := svc.db.GetRunnerUsers(&database.RunnerUsersInput{
		RunnerName: flags.runnerName,
		UserID:     flags.userid,
	})
	if err != nil {
		log.WithFields(logrus.Fields{
			"action": "RunDDBRunnerList",
			"error":  err.Error(),
		}).Error("error getting runner users")
		return err
	}
	for _, user := range res {
		logrus.WithFields(logrus.Fields{
			"action":      "RunDDBRunnerList",
			"userid":      user.UserID,
			"flags":       user.Flags,
			"newFlagsBin": fmt.Sprintf("%08b", user.Flags),
			"favorites":   Has(user.Flags, database.F_favorites),
			"followers":   Has(user.Flags, database.F_followers),
			"friends":     Has(user.Flags, database.F_friends),
			"timeline":    Has(user.Flags, database.F_timeline),
		}).Info("flags")
	}
	return nil
}
