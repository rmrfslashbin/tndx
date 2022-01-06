package params

import (
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"
	"time"

	"github.com/rmrfslashbin/tndx/pkg/database"
	"github.com/rmrfslashbin/tndx/pkg/service"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

type Outputs struct {
	Favorites *database.FavoritesItem `json:"favorites" yaml:"favorites"`
	Followers *database.FollowersItem `json:"followers" yaml:"followers"`
	Friends   *database.FriendsItem   `json:"friends" yaml:"friends"`
	Timeline  *database.TweetsItem    `json:"timeline" yaml:"timeline"`
}

func runDDBPramsGet() error {
	if flags.userid == 0 && flags.screenname != "" {
		user, _, err := svc.twitter.GetUser(&service.QueryParams{ScreenName: flags.screenname, UserID: flags.userid})
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

	outputs := &Outputs{}

	if resp, err := svc.db.GetFavoritesConfig(flags.userid); err != nil {
		log.WithFields(logrus.Fields{
			"action": "runDDBPramsGet::GetFavoritesConfig",
			"error":  err.Error(),
			"userid": flags.userid,
		}).Error("error getting favorites config")
		return err
	} else {
		outputs.Favorites = resp
		outputs.Favorites.LastUpdateTimestamp = time.UnixMilli(resp.LastUpdate)
	}

	if resp, err := svc.db.GetFollowersConfig(flags.userid); err != nil {
		log.WithFields(logrus.Fields{
			"action": "runDDBPramsGet::GetFollowersConfig",
			"error":  err.Error(),
			"userid": flags.userid,
		}).Error("error getting followers config")
		return err
	} else {
		outputs.Followers = resp
		outputs.Followers.LastUpdateTimestamp = time.UnixMilli(resp.LastUpdate)
	}

	if resp, err := svc.db.GetFriendsConfig(flags.userid); err != nil {
		log.WithFields(logrus.Fields{
			"action": "runDDBPramsGet::GetFriendsConfig",
			"error":  err.Error(),
			"userid": flags.userid,
		}).Error("error getting friends config")
		return err
	} else {
		outputs.Friends = resp
		outputs.Friends.LastUpdateTimestamp = time.UnixMilli(resp.LastUpdate)
	}

	if resp, err := svc.db.GetTimelineConfig(flags.userid); err != nil {
		log.WithFields(logrus.Fields{
			"action": "runDDBPramsGet::GetTimelineConfig",
			"error":  err.Error(),
			"userid": flags.userid,
		}).Error("error getting timeline config")
		return err
	} else {
		outputs.Timeline = resp
		outputs.Timeline.LastUpdateTimestamp = time.UnixMilli(resp.LastUpdate)
	}

	if flags.json {
		if data, err := json.Marshal(outputs); err != nil {
			log.WithFields(logrus.Fields{
				"error": err,
			}).Error("error marshalling users to json")
			return err
		} else {
			os.Stdout.Write(data)
		}
	} else if flags.yaml {
		if data, err := yaml.Marshal(outputs); err != nil {
			log.WithFields(logrus.Fields{
				"error": err,
			}).Error("error marshalling users to yaml")
			return err
		} else {
			os.Stdout.Write(data)
		}
	} else {
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.Debug)
		fmt.Fprintln(w, "UserID\tDomain\tLastUpdate\tSinceID\tMaxID")
		fmt.Fprintf(w, "%d\tfavorites\t%s\t%d\t%d\n", flags.userid, outputs.Favorites.LastUpdateTimestamp.Format(time.RFC3339), outputs.Favorites.SinceID, outputs.Favorites.MaxID)
		fmt.Fprintf(w, "%d\ttimeline\t%s\t%d\t%d\n", flags.userid, outputs.Timeline.LastUpdateTimestamp.Format(time.RFC3339), outputs.Timeline.SinceID, outputs.Timeline.MaxID)
		fmt.Fprintf(w, "\n")
		fmt.Fprintln(w, "UserID\tDomain\tLastUpdate\tPreviousCursor\tNextCursor")
		fmt.Fprintf(w, "%d\tfollowers\t%s\t%d\t%d\n", flags.userid, outputs.Followers.LastUpdateTimestamp.Format(time.RFC3339), outputs.Followers.PreviousCursor, outputs.Followers.NextCursor)
		fmt.Fprintf(w, "%d\tfriends\t%s\t%d\t%d\n", flags.userid, outputs.Friends.LastUpdateTimestamp.Format(time.RFC3339), outputs.Friends.PreviousCursor, outputs.Friends.NextCursor)
		w.Flush()
		fmt.Println()
	}

	return nil
}
