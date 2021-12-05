package subcmds

import (
	"github.com/sirupsen/logrus"
)

func RunEntitiesCmd() error {
	if err := svc.queue.SendMessage("1466090637544112141", "https://pbs.twimg.com/media/FFiZqstXwA0FsEW.jpg"); err != nil {
		log.WithFields(logrus.Fields{
			"action": "RunEntitiesCmd::SendMessage",
			"error":  err.Error(),
		}).Error("error sending message")
		return err
	}
	return nil
}
