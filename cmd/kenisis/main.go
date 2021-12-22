package main

import (
	"io/ioutil"
	"os"
	"path"

	"github.com/rmrfslashbin/tndx/pkg/kenisis"
	"github.com/sirupsen/logrus"
)

const (
	region         = "us-east-2"
	deliveryStream = "PUT-S3-12bor"
	jsonPath       = "/Users/rmrfslashbin/tndx2/timeline/16020064"
)

func main() {
	log := logrus.New()
	log.SetLevel(logrus.InfoLevel)
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	k := kenisis.NewFirehose(
		kenisis.SetRegion(region),
		kenisis.SetLogger(log),
		kenisis.SetDeliveryStream(deliveryStream),
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
		if opt, err := k.PutRecord(byteValue); err != nil {
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
