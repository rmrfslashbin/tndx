package ddb

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/rmrfslashbin/tndx/pkg/service"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

type Friends struct {
	UserID  *int64   `json:"user_id" yaml:"user_id"`
	Friends *[]int64 `json:"friends" yaml:"friends"`
	Count   int      `json:"count" yaml:"count"`
}

type UsersByFriend struct {
	FriendID *int64   `json:"friend" yaml:"friend"`
	UserIDs  *[]int64 `json:"user_ids" yaml:"user_ids"`
	Count    int      `json:"count" yaml:"count"`
}

func runDDBFriends() error {
	if flags.friendid != 0 {
		return byFriendId()
	} else {
		return friendsByUserId()
	}
}

func byFriendId() error {
	res, err := svc.db.GetFriendsByFriendId(flags.friendid)
	if err != nil {
		log.WithFields(logrus.Fields{
			"action": "runDDBFriends::GetFriendsByFriendId",
			"error":  err.Error(),
			"userid": flags.userid,
		}).Error("error getting users for friend")
		return err
	}
	users := make([]int64, len(res))
	for i, v := range res {
		users[i] = v.FriendID
	}

	results := &UsersByFriend{
		FriendID: &flags.friendid,
		UserIDs:  &users,
		Count:    len(users),
	}

	if flags.json {
		if data, err := json.Marshal(results); err != nil {
			log.WithFields(logrus.Fields{
				"error":  err,
				"action": "runDDBFriends::json.Marshal",
			}).Error("error marshalling friends to json")
			return err
		} else {
			os.Stdout.Write(data)
		}
	} else if flags.yaml {
		if data, err := yaml.Marshal(results); err != nil {
			log.WithFields(logrus.Fields{
				"error":  err,
				"action": "runDDBFriends::yaml.Marshal",
			}).Error("error marshalling friends to yaml")
			return err
		} else {
			os.Stdout.Write(data)
		}
	} else {
		fmt.Printf("Found %d users for friend: %d\n", results.Count, *results.FriendID)
		for _, v := range res {
			fmt.Printf("%d\n", v.UserID)
		}
	}

	return nil
}

func friendsByUserId() error {
	if flags.userid == 0 && flags.screenname != "" {
		user, _, err := svc.twitter.GetUser(&service.QueryParams{ScreenName: flags.screenname, UserID: flags.userid})
		if err != nil {
			log.WithFields(logrus.Fields{
				"action": "runDDBFriends::GetUser",
				"error":  err.Error(),
			}).Error("error getting user")
			return err
		}
		// Set the userid.
		flags.userid = user.ID
	}

	if flags.userid == 0 {
		log.WithFields(logrus.Fields{
			"action": "runDDBFriends::GetUser",
			"error":  "no userid or screenname provided/could not be resolved",
		}).Fatal("error getting user")
	}

	res, err := svc.db.GetFriendsByUserId(flags.userid)
	if err != nil {
		log.WithFields(logrus.Fields{
			"action": "runDDBFriends::GetFriends",
			"error":  err.Error(),
			"userid": flags.userid,
		}).Error("error getting friends for user")
		return err
	}
	friends := make([]int64, len(res))
	for i, v := range res {
		friends[i] = v.FriendID
	}

	results := &Friends{
		UserID:  &flags.userid,
		Friends: &friends,
		Count:   len(friends),
	}

	if flags.json {
		if data, err := json.Marshal(results); err != nil {
			log.WithFields(logrus.Fields{
				"error":  err,
				"action": "runDDBFriends::json.Marshal",
			}).Error("error marshalling friends to json")
			return err
		} else {
			os.Stdout.Write(data)
		}
	} else if flags.yaml {
		if data, err := yaml.Marshal(results); err != nil {
			log.WithFields(logrus.Fields{
				"error":  err,
				"action": "runDDBFriends::yaml.Marshal",
			}).Error("error marshalling friends to yaml")
			return err
		} else {
			os.Stdout.Write(data)
		}
	} else {
		fmt.Printf("Found %d friends for user: %d\n", results.Count, *results.UserID)
		for _, v := range res {
			fmt.Printf("%d\n", v.FriendID)
		}
	}
	return nil
}
