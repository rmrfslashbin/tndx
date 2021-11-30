package subcmds

import (
	"encoding/json"
	"path"
	"strconv"

	"github.com/rmrfslashbin/tndx/pkg/database"
	"github.com/rmrfslashbin/tndx/pkg/service"
	"github.com/sirupsen/logrus"
)

func RunFollowersCmd() error {

	if flags.userid == 0 {
		user, _, err := svc.twitterClient.GetUser(&service.QueryParams{ScreenName: flags.screenname, UserID: flags.userid})
		if err != nil {
			log.WithFields(logrus.Fields{
				"action": "RunFollowersCmd::GetUser",
				"error":  err.Error(),
			}).Error("error getting user")
			return err
		}
		// Set the userid.
		flags.userid = user.ID
	}

	var cursor int64 = 0
	if flags.nextCursor == 0 {
		followersConfig, err := svc.db.GetFollowersConfig(flags.userid)
		if err != nil {
			log.WithFields(logrus.Fields{
				"action": "RunFollowersCmd::GetFollowersConfig",
				"error":  err.Error(),
			}).Error("error getting follower config")
			return err
		}
		log.WithFields(logrus.Fields{
			"action":         "RunFollowersCmd::GetFollowersConfig",
			"userid":         flags.userid,
			"nextCursor":     followersConfig.NextCursor,
			"previousCursor": followersConfig.PreviousCursor,
		}).Debug("got followers config")

		if flags.backwards {
			cursor = followersConfig.PreviousCursor
		} else {
			cursor = followersConfig.NextCursor
		}
	} else {
		cursor = flags.nextCursor
	}

	log.WithFields(logrus.Fields{
		"action":    "RunFollowersCmd::Setup",
		"userid":    flags.userid,
		"backwards": flags.backwards,
		"cursor":    cursor,
	}).Debug("setting up followers")

	followers, resp, err := svc.twitterClient.GetUserFollowers(
		&service.QueryParams{
			Count:  200,
			UserID: flags.userid,
			Cursor: cursor,
		},
	)
	if err != nil {
		log.WithFields(logrus.Fields{
			"action":   "RunFollowersCmd::GetUserFollowers",
			"response": resp.Status,
		}).Error("error getting user's followers")
		return err
	}
	//spew.Dump(followers)

	if err := svc.db.PutFollowersConfig(
		&database.CursoredTweetConfigQuery{
			UserID:         flags.userid,
			NextCursor:     followers.NextCursor,
			PreviousCursor: followers.PreviousCursor,
		},
	); err != nil {
		logrus.WithFields(logrus.Fields{
			"action":         "RunFollowersCmd::PutFollowersConfig",
			"error":          err.Error(),
			"userid":         flags.userid,
			"nextCursor":     followers.NextCursor,
			"previousCursor": followers.PreviousCursor,
		}).Error("error putting followers config")
		return err
	}

	// Save the users.
	for f := range followers.Users {
		if data, err := json.MarshalIndent(followers.Users[f], "", "  "); err == nil {
			if err := svc.storage.Put(path.Join(strconv.FormatInt(flags.userid, 10), "followers", followers.Users[f].IDStr), data); err != nil {
				logrus.WithFields(logrus.Fields{
					"action":     "RunFollowersCmd::Put",
					"error":      err.Error(),
					"userid":     flags.userid,
					"followerId": followers.Users[f].ID,
				}).Error("error putting followers")
				return err
			}
		}
	}
	logrus.WithFields(logrus.Fields{
		"action":         "RunFollowersCmd::Done!",
		"userid":         flags.userid,
		"nextCursor":     followers.NextCursor,
		"previousCursor": followers.PreviousCursor,
		"count":          len(followers.Users),
	}).Info("finished getting followers")

	return nil
}
