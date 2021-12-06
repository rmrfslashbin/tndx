package subcmds

import (
	"github.com/rmrfslashbin/tndx/pkg/service"
	"github.com/sirupsen/logrus"
)

func RunDDBRunnerDel() error {
	if flags.userid == 0 {
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

	err := svc.db.DeleteRunnerUser(flags.runnerName, flags.userid)
	if err != nil {
		log.WithFields(logrus.Fields{
			"action": "RunDDBRunnerSet",
			"error":  err.Error(),
		}).Error("error deleteing runner user")
		return err
	}

	return nil
}
