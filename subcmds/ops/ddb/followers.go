package ddb

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/rmrfslashbin/tndx/pkg/service"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

type Followers struct {
	UserID    *int64   `json:"user_id" yaml:"user_id"`
	Followers *[]int64 `json:"followers" yaml:"followers"`
	Count     int      `json:"count" yaml:"count"`
}

type UserByFollower struct {
	UserID   *[]int64 `json:"user_ids" yaml:"user_ids"`
	Follower *int64   `json:"follower" yaml:"follower"`
	Count    int      `json:"count" yaml:"count"`
}

func runDDBFollowers() error {
	if flags.followid != 0 {
		return byFolloweId()
	} else {
		return followersByUserId()
	}
}

func byFolloweId() error {
	res, err := svc.db.GetFollowersByFollowId(flags.followid)
	if err != nil {
		log.WithFields(logrus.Fields{
			"action": "runDDBFollowers::GetFollowersByFollowId",
			"error":  err.Error(),
			"userid": flags.userid,
		}).Error("error getting users for follower")
		return err
	}
	followers := make([]int64, len(res))
	for i, v := range res {
		followers[i] = v.FollowerID
	}

	results := &UserByFollower{
		UserID:   &followers,
		Follower: &flags.followid,
		Count:    len(followers),
	}

	if flags.json {
		if data, err := json.Marshal(results); err != nil {
			log.WithFields(logrus.Fields{
				"error":  err,
				"action": "runDDBFollowers::json.Marshal",
			}).Error("error marshalling users to json")
			return err
		} else {
			os.Stdout.Write(data)
		}
	} else if flags.yaml {
		if data, err := yaml.Marshal(results); err != nil {
			log.WithFields(logrus.Fields{
				"error":  err,
				"action": "runDDBFollowers::yaml.Marshal",
			}).Error("error marshalling users to yaml")
			return err
		} else {
			os.Stdout.Write(data)
		}
	} else {
		fmt.Printf("Found %d users for follower: %d\n", results.Count, *results.Follower)
		for _, v := range res {
			fmt.Printf("%d\n", v.FollowerID)
		}
	}

	return nil
}

func followersByUserId() error {
	if flags.userid == 0 && flags.screenname != "" {
		user, _, err := svc.twitter.GetUser(&service.QueryParams{ScreenName: flags.screenname, UserID: flags.userid})
		if err != nil {
			log.WithFields(logrus.Fields{
				"action": "runDDBFollowers::GetUser",
				"error":  err.Error(),
			}).Error("error getting user")
			return err
		}
		// Set the userid.
		flags.userid = user.ID
	}

	if flags.userid == 0 {
		log.WithFields(logrus.Fields{
			"action": "runDDBFollowers::GetUser",
			"error":  "no userid or screenname provided/could not be resolved",
		}).Fatal("error getting user")
	}

	res, err := svc.db.GetFollowersByUserId(flags.userid)
	if err != nil {
		log.WithFields(logrus.Fields{
			"action": "runDDBFollowers::GetFollowers",
			"error":  err.Error(),
			"userid": flags.userid,
		}).Error("error getting followers for user")
		return err
	}
	followers := make([]int64, len(res))
	for i, v := range res {
		followers[i] = v.FollowerID
	}

	results := &Followers{
		UserID:    &flags.userid,
		Followers: &followers,
		Count:     len(followers),
	}

	if flags.json {
		if data, err := json.Marshal(results); err != nil {
			log.WithFields(logrus.Fields{
				"error":  err,
				"action": "runDDBFollowers::json.Marshal",
			}).Error("error marshalling followers to json")
			return err
		} else {
			os.Stdout.Write(data)
		}
	} else if flags.yaml {
		if data, err := yaml.Marshal(results); err != nil {
			log.WithFields(logrus.Fields{
				"error":  err,
				"action": "runDDBFollowers::yaml.Marshal",
			}).Error("error marshalling followers to yaml")
			return err
		} else {
			os.Stdout.Write(data)
		}
	} else {
		fmt.Printf("Found %d followers for user: %d\n", results.Count, *results.UserID)
		for _, v := range res {
			fmt.Printf("%d\n", v.FollowerID)
		}
	}

	return nil
}
