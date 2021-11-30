package subcmds

import (
	"github.com/sirupsen/logrus"
)

func RunEntitiesCmd() error {
	if err := svc.queue.ReceiveMessage(); err != nil {
		log.WithFields(logrus.Fields{
			"action": "RunEntitiesCmd::ReceiveMessage",
			"error":  err.Error(),
		}).Error("error receiving message")
		return err
	}
	return nil
}
