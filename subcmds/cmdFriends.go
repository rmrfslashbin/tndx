package subcmds

import (
	"encoding/json"
	"path"
	"strconv"

	"github.com/rmrfslashbin/tndx/pkg/database"
	"github.com/rmrfslashbin/tndx/pkg/service"
	"github.com/sirupsen/logrus"
)

func RunFriendsCmd() error {

	if flags.userid == 0 {
		user, _, err := svc.twitterClient.GetUser(&service.QueryParams{ScreenName: flags.screenname, UserID: flags.userid})
		if err != nil {
			log.WithFields(logrus.Fields{
				"action": "RunFriendsCmd::GetUser",
				"error":  err.Error(),
			}).Error("error getting user")
			return err
		}
		// Set the userid.
		flags.userid = user.ID
	}

	var cursor int64 = 0
	if flags.nextCursor == 0 {
		friendsConfig, err := svc.db.GetFriendsConfig(flags.userid)
		if err != nil {
			log.WithFields(logrus.Fields{
				"action": "RunFriendsCmd::GetFriendsConfig",
				"error":  err.Error(),
			}).Error("error getting friends config")
			return err
		}
		log.WithFields(logrus.Fields{
			"action":     "RunFriendsCmd::GetFriendsConfig",
			"userid":     flags.userid,
			"nextCursor": friendsConfig.NextCursor,
		}).Debug("got friends config")

		if flags.backwards {
			cursor = friendsConfig.PreviousCursor
		} else {
			cursor = friendsConfig.NextCursor
		}
	} else {
		cursor = flags.nextCursor
	}

	log.WithFields(logrus.Fields{
		"action":    "RunFriendsCmd::Setup",
		"userid":    flags.userid,
		"backwards": flags.backwards,
		"cursor":    cursor,
	}).Debug("setting up friends")

	friends, resp, err := svc.twitterClient.GetUserFriends(
		&service.QueryParams{
			Count:  200,
			UserID: flags.userid,
			Cursor: cursor,
		},
	)
	if err != nil {
		log.WithFields(logrus.Fields{
			"action":   "RunFriendsCmd::GetUserFriends",
			"response": resp.Status,
		}).Error("error getting user's friends")
		return err
	}
	//spew.Dump(friends)

	if err := svc.db.PutFriendsConfig(
		&database.CursoredTweetConfigQuery{
			UserID:         flags.userid,
			NextCursor:     friends.NextCursor,
			PreviousCursor: friends.PreviousCursor,
		},
	); err != nil {
		logrus.WithFields(logrus.Fields{
			"action":         "RunFriendsCmd::PutFriendsConfig",
			"error":          err.Error(),
			"userid":         flags.userid,
			"nextCursor":     friends.NextCursor,
			"previousCursor": friends.PreviousCursor,
		}).Error("error putting friends config")
		return err
	}

	// Save the users.
	for f := range friends.Users {
		if data, err := json.Marshal(friends.Users[f]); err == nil {
			if err := svc.storage.Put(path.Join("friends", strconv.FormatInt(flags.userid, 10), friends.Users[f].IDStr+".json"), data); err != nil {
				logrus.WithFields(logrus.Fields{
					"action":   "RunFriendsCmd::Put",
					"error":    err.Error(),
					"userid":   flags.userid,
					"friendId": friends.Users[f].ID,
				}).Error("error putting friends")
				return err
			}
		}
	}
	logrus.WithFields(logrus.Fields{
		"action":         "RunFriendsCmd::Done!",
		"userid":         flags.userid,
		"nextCursor":     friends.NextCursor,
		"previousCursor": friends.PreviousCursor,
		"count":          len(friends.Users),
	}).Info("finished getting friends")

	return nil
}
