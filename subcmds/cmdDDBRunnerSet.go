package subcmds

import (
	"fmt"

	"github.com/rmrfslashbin/tndx/pkg/database"
	"github.com/rmrfslashbin/tndx/pkg/service"
	"github.com/sirupsen/logrus"
)

func RunDDBRunnerSet() error {
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

	var newFlags database.Bits
	if flags.favorites {
		newFlags = Set(newFlags, database.F_favorites)
	}

	if flags.followers {
		newFlags = Set(newFlags, database.F_followers)
	}

	if flags.friends {
		newFlags = Set(newFlags, database.F_friends)
	}

	if flags.timeline {
		newFlags = Set(newFlags, database.F_timeline)
	}

	if flags.user {
		newFlags = Set(newFlags, database.F_user)
	}

	if flags.all {
		newFlags = Set(newFlags, database.F_favorites)
		newFlags = Set(newFlags, database.F_followers)
		newFlags = Set(newFlags, database.F_friends)
		newFlags = Set(newFlags, database.F_timeline)
		newFlags = Set(newFlags, database.F_user)
	}

	logrus.WithFields(logrus.Fields{
		"action":      "RunDDBRunnerSet",
		"userid":      flags.userid,
		"screenname":  flags.screenname,
		"newFlags":    newFlags,
		"newFlagsBin": fmt.Sprintf("%08b", newFlags),
		"favorites":   Has(newFlags, database.F_favorites),
		"followers":   Has(newFlags, database.F_followers),
		"friends":     Has(newFlags, database.F_friends),
		"timeline":    Has(newFlags, database.F_timeline),
		"user":        Has(newFlags, database.F_user),
	}).Info("setting flags")

	if err := svc.db.PutRunnerFlags(flags.runnerName, flags.userid, newFlags); err != nil {
		log.WithFields(logrus.Fields{
			"action": "RunDDBRunnerSet",
			"error":  err.Error(),
		}).Error("error setting flags")
		return err
	}

	return nil
}

func Set(b, flag database.Bits) database.Bits    { return b | flag }
func Clear(b, flag database.Bits) database.Bits  { return b &^ flag }
func Toggle(b, flag database.Bits) database.Bits { return b ^ flag }
func Has(b, flag database.Bits) bool             { return b&flag != 0 }
