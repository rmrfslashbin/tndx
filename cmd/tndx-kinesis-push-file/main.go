package main

import (
	"encoding/json"
	"io"
	"os"
	"path"
	"time"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/rmrfslashbin/tndx/pkg/kinesis"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

const (
	region         = "us-east-2"
	deliveryStream = "tndx-rmrfslashbin-DeliveryStreamTweets-DPg7a3ClhQ7u"
)

var (
	RootCmd = &cobra.Command{
		Version: "v2021.12.20-00",
		Run: func(cmd *cobra.Command, args []string) {
			run()
		},
	}
	jsonFile string
	log      *logrus.Logger
	k        *kinesis.Config
)

func init() {
	log = logrus.New()
	log.SetLevel(logrus.InfoLevel)
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	k = kinesis.NewFirehose(
		kinesis.SetRegion(region),
		kinesis.SetLogger(log),
		kinesis.SetDeliveryStream(deliveryStream),
	)

	RootCmd.Flags().StringVarP(&jsonFile, "jsonfile", "j", "", "path to json file")
	RootCmd.MarkFlagRequired("jsonfile")
}

func run() {
	fqpn := path.Clean(jsonFile)

	jsonFile, err := os.Open(fqpn)
	if err != nil {
		log.WithFields(logrus.Fields{
			"error": err,
			"file":  fqpn,
		}).Fatal("failed opening file")
	}
	defer jsonFile.Close()
	log.WithFields(logrus.Fields{
		"file": fqpn,
	}).Info("opened file")

	dec := json.NewDecoder(jsonFile)
	count := 0
	subcount := 0
	for {
		count++
		subcount++
		if subcount >= 500 {
			log.Info("Taking a break...")
			time.Sleep(200 * time.Second)
			subcount = 0
		}
		var tweet twitter.Tweet
		if err := dec.Decode(&tweet); err != nil {
			if err == io.EOF {
				break
			}
			if err != nil {
				log.WithFields(logrus.Fields{
					"error": err,
					"file":  fqpn,
				}).Fatal("failed unmarshalling file")
			}
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
				"recordId": *opt.RecordId,
			}).Info("put record")
		}
	}
	log.WithFields(logrus.Fields{
		"file":  fqpn,
		"count": count,
	}).Info("processed file")
}

func main() {
	RootCmd.Execute()
}
