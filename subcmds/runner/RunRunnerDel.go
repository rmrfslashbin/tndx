package runner

import (
	"github.com/rmrfslashbin/tndx/pkg/database"
	"github.com/rmrfslashbin/tndx/pkg/service"
	"github.com/sirupsen/logrus"
)

func RunRunnerDel() error {
	if flags.userid == 0 && flags.screenname != "" {
		user, _, err := svc.twitterClient.GetUser(&service.QueryParams{ScreenName: flags.screenname, UserID: flags.userid})
		if err != nil {
			log.WithFields(logrus.Fields{
				"action": "RunRunnerDel::GetUser",
				"error":  err.Error(),
			}).Error("error getting user")
			return err
		}
		// Set the userid.
		flags.userid = user.ID
	}

	if flags.userid == 0 {
		log.WithFields(logrus.Fields{
			"action": "RunRunnerDel::GetUser",
			"error":  "no userid or screenname provided/could not be resolved",
		}).Fatal("error getting user")
	}

	if err := svc.db.DeleteRunnerUser(&database.RunnerFlagsItem{RunnerName: flags.runner, UserID: flags.userid}); err != nil {
		log.WithFields(logrus.Fields{
			"action": "RunRunnerDel::PutRunnerFlags",
			"error":  err.Error(),
		}).Error("error putting runner flags")
		return err
	}
	log.Info("runner user deleted")
	return nil
}
