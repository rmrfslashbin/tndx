package users

import (
	"encoding/json"
	"os"

	"github.com/davecgh/go-spew/spew"
	"github.com/dghubble/go-twitter/twitter"
	"github.com/rmrfslashbin/tndx/pkg/service"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

type BasicData struct {
	CreatedAt       string `json:"created_at"`
	Description     string `json:"description"`
	FavouritesCount int    `json:"favourites_count"`
	FollowersCount  int    `json:"followers_count"`
	FriendsCount    int    `json:"friends_count"`
	Location        string `json:"location"`
	Name            string `json:"name"`
	Protected       bool   `json:"protected"`
	ScreenName      string `json:"screen_name"`
	StatusesCount   int    `json:"statuses_count"`
	URL             string `json:"url"`
	UserID          int64  `json:"user_id"`
	Verified        bool   `json:"verified"`
}

func runUsersGet() error {
	userIDs := DedupInt64Slice(flags.userids)
	screenNames := DedupStringSlice(flags.screenname)
	users, resp, err := svc.twitter.LookupUsers(&twitter.UserLookupParams{
		UserID:     userIDs,
		ScreenName: screenNames,
	})
	if err != nil {
		log.WithFields(logrus.Fields{
			"error":       err,
			"code":        resp.StatusCode,
			"status":      resp.Status,
			"userIDs":     userIDs,
			"screenNames": screenNames,
		}).Error("error looking up users")
		return err
	}
	log.WithFields(logrus.Fields{
		"count": len(users),
	}).Info("users returned")

	if flags.basic {
		basicData := make([]BasicData, len(users))
		for u := range users {
			users[u].CreatedAt, _ = service.FixTwitterTimeRFC3339(users[u].CreatedAt)
			if users[u].Status != nil {
				users[u].Status.CreatedAt, _ = service.FixTwitterTimeRFC3339(users[u].Status.CreatedAt)
			}
			basicData[u] = BasicData{
				CreatedAt:       users[u].CreatedAt,
				Description:     users[u].Description,
				FavouritesCount: users[u].FavouritesCount,
				FollowersCount:  users[u].FollowersCount,
				FriendsCount:    users[u].FriendsCount,
				Location:        users[u].Location,
				Name:            users[u].Name,
				Protected:       users[u].Protected,
				ScreenName:      users[u].ScreenName,
				StatusesCount:   users[u].StatusesCount,
				URL:             users[u].URL,
				UserID:          users[u].ID,
				Verified:        users[u].Verified,
			}
		}
		if flags.json {
			if data, err := json.Marshal(basicData); err != nil {
				log.WithFields(logrus.Fields{
					"error": err,
				}).Error("error marshalling users to json")
				return err
			} else {
				os.Stdout.Write(data)
			}
		} else if flags.yaml {
			if data, err := yaml.Marshal(basicData); err != nil {
				log.WithFields(logrus.Fields{
					"error": err,
				}).Error("error marshalling users to yaml")
				return err
			} else {
				os.Stdout.Write(data)
			}
		} else {
			spew.Dump(basicData)
		}
	} else {
		for u := range users {
			users[u].CreatedAt, _ = service.FixTwitterTime(users[u].CreatedAt)
			if users[u].Status != nil {
				users[u].Status.CreatedAt, _ = service.FixTwitterTime(users[u].Status.CreatedAt)
			}
		}
		if flags.json {
			if data, err := json.Marshal(users); err != nil {
				log.WithFields(logrus.Fields{
					"error": err,
				}).Error("error marshalling users to json")
				return err
			} else {
				os.Stdout.Write(data)
			}
		} else if flags.yaml {
			if data, err := yaml.Marshal(users); err != nil {
				log.WithFields(logrus.Fields{
					"error": err,
				}).Error("error marshalling users to yaml")
				return err
			} else {
				os.Stdout.Write(data)
			}
		} else {
			spew.Dump(users)
		}
	}

	return nil

}
