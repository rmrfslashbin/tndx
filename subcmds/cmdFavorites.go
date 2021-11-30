package subcmds

import (
	"encoding/json"
	"path"
	"strconv"

	"github.com/rmrfslashbin/tndx/pkg/database"
	"github.com/rmrfslashbin/tndx/pkg/service"
	"github.com/sirupsen/logrus"
)

func RunFavoritesCmd() error {

	if flags.userid == 0 {
		user, _, err := svc.twitterClient.GetUser(&service.QueryParams{ScreenName: flags.screenname, UserID: flags.userid})
		if err != nil {
			log.WithFields(logrus.Fields{
				"action": "RunFavoritesCmd::GetUser",
				"error":  err.Error(),
			}).Error("error getting user")
			return err
		}
		// Set the userid.
		flags.userid = user.ID
	}

	var maxID int64
	var sinceID int64
	if flags.maxid != 0 || flags.sinceid != 0 {
		maxID = flags.maxid
		sinceID = flags.sinceid
	} else {
		favConfig, err := svc.db.GetFavoritesConfig(flags.userid)
		if err != nil {
			log.WithFields(logrus.Fields{
				"action": "RunFavoritesCmd::GetFavoritesConfig",
				"error":  err.Error(),
			}).Error("error getting favorites config")
			return err
		}

		if flags.backwards {
			maxID = favConfig.MaxID
		} else {
			sinceID = favConfig.SinceID
		}
	}

	log.WithFields(logrus.Fields{
		"action":    "RunFavoritesCmd::Setup",
		"userid":    flags.userid,
		"maxid":     maxID,
		"sinceid":   sinceID,
		"backwards": flags.backwards,
	}).Debug("setting up favorites")

	tweets, resp, err := svc.twitterClient.GetUserFavorites(
		&service.QueryParams{
			Count:   200,
			MaxID:   maxID,
			SinceID: sinceID,
			UserID:  flags.userid,
		},
	)
	if err != nil {
		log.WithFields(logrus.Fields{
			"action":   "RunFavoritesCmd::GetUserFriends",
			"response": resp.Status,
		}).Error("error getting user's favorites")
		return err
	}

	// upperID and lowerID are used to keep track of the max and min tweet IDs
	var upperID int64
	var lowerID int64

	// Loop through all the tweets.
	for t := range tweets {
		if data, err := json.MarshalIndent(tweets[t], "", "  "); err == nil {
			if err := svc.storage.Put(path.Join(strconv.FormatInt(flags.userid, 10), "favorites", tweets[t].IDStr), data); err != nil {
				logrus.WithFields(logrus.Fields{
					"action":  "RunFavoritesCmd::Put",
					"error":   err.Error(),
					"userid":  flags.userid,
					"tweetId": tweets[t].ID,
				}).Error("error putting friends")
				return err
			}
		}

		// Calculate the max and min tweet IDs.
		if tweets[t].ID > upperID {
			upperID = tweets[t].ID
		}
		if tweets[t].ID < lowerID {
			lowerID = tweets[t].ID
		}
		if lowerID == 0 {
			lowerID = upperID
		}
	}

	if upperID > 0 {
		if err := svc.db.PutFavoritesConfig(
			&database.TweetConfigQuery{
				UserID:  flags.userid,
				SinceID: lowerID,
				MaxID:   upperID,
			},
		); err != nil {
			logrus.WithFields(logrus.Fields{
				"action":       "RunFavoritesCmd::PutFriendsConfig",
				"error":        err.Error(),
				"userid":       flags.userid,
				"MaxUpperID":   upperID,
				"SinceLowerID": lowerID,
			}).Error("error putting favorites config")
			return err
		}
	}

	logrus.WithFields(logrus.Fields{
		"action":  "RunFavoritesCmd::Done!",
		"userid":  flags.userid,
		"upperID": upperID,
		"lowerID": lowerID,
		"count":   len(tweets),
	}).Info("finished getting favorites")

	return nil
}
