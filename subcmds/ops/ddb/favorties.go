package ddb

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/rmrfslashbin/tndx/pkg/service"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

type Favorites struct {
	UserID    *int64   `json:"user_id" yaml:"user_id"`
	Favorites *[]int64 `json:"favorites" yaml:"favorites"`
	Count     int      `json:"count" yaml:"count"`
}

type UsersByFavorite struct {
	FavoriteID *int64   `json:"favorite" yaml:"favorite"`
	UserIDs    *[]int64 `json:"user_ids" yaml:"user_ids"`
	Count      int      `json:"count" yaml:"count"`
}

func runDDBFavorites() error {
	if flags.tweetid != 0 {
		return byTweetId()
	} else {
		return favoritesByUserId()
	}
}

func byTweetId() error {
	res, err := svc.db.GetFavoritesByTweetId(flags.tweetid)
	if err != nil {
		log.WithFields(logrus.Fields{
			"action":  "runDDBFavorites::GetFavoritesByTweetId",
			"error":   err.Error(),
			"tweetid": flags.tweetid,
		}).Error("error getting users for favorite tweet")
		return err
	}

	users := make([]int64, len(res))
	for i, v := range res {
		users[i] = v.UserID
	}
	results := &UsersByFavorite{
		UserIDs:    &users,
		FavoriteID: &flags.tweetid,
		Count:      len(users),
	}

	if flags.json {
		if data, err := json.Marshal(results); err != nil {
			log.WithFields(logrus.Fields{
				"error":  err,
				"action": "runDDBFavorites::json.Marshal",
			}).Error("error marshalling favorites to json")
			return err
		} else {
			os.Stdout.Write(data)
		}
	} else if flags.yaml {
		if data, err := yaml.Marshal(results); err != nil {
			log.WithFields(logrus.Fields{
				"error":  err,
				"action": "runDDBFavorites::yaml.Marshal",
			}).Error("error marshalling favorites to yaml")
			return err
		} else {
			os.Stdout.Write(data)
		}
	} else {
		fmt.Printf("Found %d favorites for tewwt: %d\n", results.Count, *results.FavoriteID)
		for _, v := range res {
			fmt.Printf("%d\n", v.UserID)
		}
	}

	return nil
}

func favoritesByUserId() error {
	if flags.userid == 0 && flags.screenname != "" {
		user, _, err := svc.twitter.GetUser(&service.QueryParams{ScreenName: flags.screenname, UserID: flags.userid})
		if err != nil {
			log.WithFields(logrus.Fields{
				"action": "runDDBFavorites::GetUser",
				"error":  err.Error(),
			}).Error("error getting user")
			return err
		}
		// Set the userid.
		flags.userid = user.ID
	}

	if flags.userid == 0 {
		return errors.New("no userid or screenname provided/could not be resolved")
	}

	res, err := svc.db.GetFavoritesByUserId(flags.userid)
	if err != nil {
		log.WithFields(logrus.Fields{
			"action": "runDDBFavorites::GetFavorites",
			"error":  err.Error(),
			"userid": flags.userid,
		}).Error("error getting favorites for user")
		return err
	}
	favorites := make([]int64, len(res))
	for i, v := range res {
		favorites[i] = v.TweetID
	}

	results := &Favorites{
		UserID:    &flags.userid,
		Favorites: &favorites,
		Count:     len(favorites),
	}

	if flags.json {
		if data, err := json.Marshal(results); err != nil {
			log.WithFields(logrus.Fields{
				"error":  err,
				"action": "runDDBFavorites::json.Marshal",
			}).Error("error marshalling favorites to json")
			return err
		} else {
			os.Stdout.Write(data)
		}
	} else if flags.yaml {
		if data, err := yaml.Marshal(results); err != nil {
			log.WithFields(logrus.Fields{
				"error":  err,
				"action": "runDDBFavorites::yaml.Marshal",
			}).Error("error marshalling favorites to yaml")
			return err
		} else {
			os.Stdout.Write(data)
		}
	} else {
		fmt.Printf("Found %d favorites for user: %d\n", results.Count, *results.UserID)
		for _, v := range res {
			fmt.Printf("%d\n", v.TweetID)
		}
	}

	return nil
}
