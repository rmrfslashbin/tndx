package subcmds

import (
	"encoding/json"
	"path"

	"github.com/davecgh/go-spew/spew"
	"github.com/rmrfslashbin/tndx/pkg/service"
	"github.com/sirupsen/logrus"
)

func RunUserCmd() error {
	user, _, err := svc.twitterClient.GetUser(&service.QueryParams{ScreenName: flags.screenname, UserID: flags.userid})
	if err != nil {
		log.WithFields(logrus.Fields{
			"action": "RunFriendsCmd::GetUser",
			"error":  err.Error(),
		}).Error("error getting user.")
		return err
	}

	spew.Dump(user)
	if data, err := json.Marshal(user); err == nil {
		svc.storage.Put(path.Join("users", user.IDStr+".json"), data)
	}
	return nil
}
