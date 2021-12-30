package main

import (
	"bufio"
	"encoding/base64"
	"encoding/json"
	"os"
	"path"

	"github.com/rmrfslashbin/tndx/pkg/kenisis"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type Record struct {
	AttemptsMade           int    `json:"attemptsMade"`
	ArrivalTimestamp       int64  `json:"arrivalTimestamp"`
	ErrorCode              string `json:"errorCode"`
	ErrorMessage           string `json:"errorMessage"`
	AttemptEndingTimestamp int64  `json:"attemptEndingTimestamp"`
	RawData                string `json:"rawData"`
}

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
	k        *kenisis.Config
)

func init() {
	log = logrus.New()
	log.SetLevel(logrus.InfoLevel)
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	k = kenisis.NewFirehose(
		kenisis.SetRegion(region),
		kenisis.SetLogger(log),
		kenisis.SetDeliveryStream(deliveryStream),
	)

	RootCmd.Flags().StringVarP(&jsonFile, "jsonfile", "j", "", "path to json file")
	RootCmd.MarkFlagRequired("jsonfile")
}

func main() {
	RootCmd.Execute()
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

	scanner := bufio.NewScanner(jsonFile)
	count := 0
	for scanner.Scan() {
		record := &Record{}
		if err := json.Unmarshal(scanner.Bytes(), record); err != nil {
			log.WithFields(logrus.Fields{
				"error": err,
				"file":  fqpn,
			}).Fatal("failed unmarshalling json")
		}
		data, err := base64.StdEncoding.DecodeString(record.RawData)
		if err != nil {
			log.WithFields(logrus.Fields{
				"error": err,
				"file":  fqpn,
			}).Fatal("failed decoding base64")
		}

		if opt, err := k.PutRecord(data); err != nil {
			log.WithFields(logrus.Fields{
				"file":  fqpn,
				"error": err,
			}).Fatal("failed putting record")
		} else {
			log.WithFields(logrus.Fields{
				"recordId": *opt.RecordId,
			}).Info("put record")
			count++
		}
	}
	if err := scanner.Err(); err != nil {
		log.WithFields(logrus.Fields{
			"error": err,
			"file":  fqpn,
		}).Fatal("failed reading file")
	}
	log.WithFields(logrus.Fields{
		"file":  fqpn,
		"count": count,
	}).Info("processed file")
}
