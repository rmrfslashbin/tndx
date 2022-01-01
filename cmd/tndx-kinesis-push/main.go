package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/rmrfslashbin/tndx/pkg/kinesis"
	"github.com/rmrfslashbin/tndx/pkg/service"
	"github.com/sirupsen/logrus"
)

const (
	region         = "us-east-2"
	deliveryStream = "tndx-rmrfslashbin-DeliveryStreamTweets-DPg7a3ClhQ7u"
	jsonPath       = "/Users/rmrfslashbin/tndx2/timeline/16020064"
)

func main() {
	log := logrus.New()
	log.SetLevel(logrus.InfoLevel)
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	k := kinesis.NewFirehose(
		kinesis.SetRegion(region),
		kinesis.SetLogger(log),
		kinesis.SetDeliveryStream(deliveryStream),
	)

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
		if err := json.Unmarshal(byteValue, tweet); err != nil {
			log.WithFields(logrus.Fields{
				"error": err,
				"file":  fqpn,
			}).Fatal("failed unmarshalling file")
		}

		tweet.CreatedAt, _ = service.FixTwitterTime(tweet.CreatedAt)
		tweet.User.CreatedAt, _ = service.FixTwitterTime(tweet.User.CreatedAt)
		if tweet.RetweetedStatus != nil {
			tweet.RetweetedStatus.CreatedAt, _ = service.FixTwitterTime(tweet.RetweetedStatus.CreatedAt)
			tweet.RetweetedStatus.User.CreatedAt, _ = service.FixTwitterTime(tweet.RetweetedStatus.User.CreatedAt)
		}

		rawBytes, err := json.Marshal(tweet)
		if err != nil {
			log.WithFields(logrus.Fields{
				"error": err,
				"file":  fqpn,
			}).Fatal("failed marshalling tweet to bytes")
		}

		if opt, err := k.PutRecord(rawBytes); err != nil {
			log.WithFields(logrus.Fields{
				"file":  fqpn,
				"error": err,
			}).Fatal("failed putting record")
		} else {

			log.WithFields(logrus.Fields{
				"file":     fqpn,
				"recordId": *opt.RecordId,
			}).Info("put record")
		}
	}
}
