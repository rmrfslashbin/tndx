package runner

import (
	"github.com/rmrfslashbin/tndx/pkg/database"
	"github.com/rmrfslashbin/tndx/pkg/service"
	"github.com/sirupsen/logrus"
)

func RunRunnerSet() error {
	if flags.userid == 0 && flags.screenname != "" {
		user, _, err := svc.twitterClient.GetUser(&service.QueryParams{ScreenName: flags.screenname, UserID: flags.userid})
		if err != nil {
			log.WithFields(logrus.Fields{
				"action": "RunRunnerSet::GetUser",
				"error":  err.Error(),
			}).Error("error getting user")
			return err
		}
		// Set the userid.
		flags.userid = user.ID
	}

	if flags.userid == 0 {
		log.WithFields(logrus.Fields{
			"action": "RunRunnerSet::GetUser",
			"error":  "no userid or screenname provided/could not be resolved",
		}).Fatal("error getting user")
	}

	var newFlags database.Bits
	if flags.favorites {
		newFlags = database.Set(newFlags, database.F_favorites)
	}

	if flags.followers {
		newFlags = database.Set(newFlags, database.F_followers)
	}

	if flags.friends {
		newFlags = database.Set(newFlags, database.F_friends)
	}

	if flags.timeline {
		newFlags = database.Set(newFlags, database.F_timeline)
	}

	if flags.user {
		newFlags = database.Set(newFlags, database.F_user)
	}

	if flags.all {
		newFlags = database.Set(newFlags, database.F_favorites)
		newFlags = database.Set(newFlags, database.F_followers)
		newFlags = database.Set(newFlags, database.F_friends)
		newFlags = database.Set(newFlags, database.F_timeline)
		newFlags = database.Set(newFlags, database.F_user)
	}

	if err := svc.db.PutRunnerFlags(&database.RunnerFlagsItem{RunnerName: flags.runner, UserID: flags.userid, Flags: newFlags}); err != nil {
		log.WithFields(logrus.Fields{
			"action": "RunRunnerSet::PutRunnerFlags",
			"error":  err.Error(),
		}).Error("error putting runner flags")
		return err
	}
	log.Info("runner flags updated for user")
	return nil
}
