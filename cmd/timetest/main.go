package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strconv"
	"time"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/sirupsen/logrus"
)

const (
	region         = "us-east-2"
	deliveryStream = "tndx-timelines"
	jsonPath       = "/Users/rmrfslashbin/tndx2/timeline/16020064"
)

func FixTwitterTime(timeStr string) (string, error) {
	const layout = "Mon Jan 2 15:04:05 -0700 2006"
	if t, err := time.Parse(layout, timeStr); err != nil {
		return "", err
	} else {
		return strconv.FormatInt(t.Unix(), 10), nil
	}
}

func main() {
	log := logrus.New()
	log.SetLevel(logrus.InfoLevel)
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	files, err := ioutil.ReadDir(jsonPath)
	if err != nil {
		panic(err)
	}
	for _, file := range files {
		fqpn := path.Clean(path.Join(jsonPath, file.Name()))

		jsonFile, err := os.Open(fqpn)
		// if we os.Open returns an error then handle it
		if err != nil {
			log.WithFields(logrus.Fields{
				"error": err,
				"file":  fqpn,
			}).Fatal("failed opening file")
		}
		defer jsonFile.Close()
		byteValue, _ := ioutil.ReadAll(jsonFile)
		tweet := &twitter.Tweet{}
		json.Unmarshal(byteValue, tweet)
		tweet.CreatedAt, _ = FixTwitterTime(tweet.CreatedAt)
		tweet.User.CreatedAt, _ = FixTwitterTime(tweet.User.CreatedAt)
		if tweet.RetweetedStatus != nil {
			tweet.RetweetedStatus.CreatedAt, _ = FixTwitterTime(tweet.RetweetedStatus.CreatedAt)
			tweet.RetweetedStatus.User.CreatedAt, _ = FixTwitterTime(tweet.RetweetedStatus.User.CreatedAt)
		}
		if bytes, err := json.Marshal(tweet); err != nil {
			log.WithFields(logrus.Fields{
				"error": err,
				"file":  fqpn,
			}).Fatal("failed marshalling file")
		} else {
			fmt.Println(string(bytes))
		}
		//spew.Dump(tweet)
		break
	}
}
