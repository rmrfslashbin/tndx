package main

import (
	"encoding/csv"
	"fmt"
	"os"

	"github.com/rmrfslashbin/tndx/pkg/comprehend"
	"github.com/sirupsen/logrus"
)

var (
	text []string
	log  *logrus.Logger
	c    *comprehend.Config
)

func init() {
	log := logrus.New()
	log.SetLevel(logrus.InfoLevel)
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	c = comprehend.New(
		comprehend.SetLogger(log),
		comprehend.SetRegion("us-east-2"),
	)
}

func main() {

	// open file
	f, err := os.Open("/Users/rmrfslashbin/text.tsv")
	if err != nil {
		log.Fatal(err)
	}

	// remember to close the file at the end of the program
	defer f.Close()

	// read csv values using csv.Reader
	csvReader := csv.NewReader(f)
	csvReader.Comma = '\t'
	data, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	text = make([]string, len(data))
	for ndx, line := range data {
		text[ndx] = line[2]
	}

	//runEntities(text)
	runSentiment(text)

}

func runEntities(text []string) error {
	resp, err := c.DetectEntities(&text)
	if err != nil {
		log.Fatal(err)
	}
	if len(resp.ErrorList) > 0 {
		for _, err := range resp.ErrorList {
			log.Error(err)
		}
	}
	for _, item := range resp.ResultList {
		if len(item.Entities) > 0 {
			fmt.Printf("Index %d\n", *item.Index)
			fmt.Printf("  Text: %s\n\n", text[*item.Index])
			for _, entity := range item.Entities {
				fmt.Printf("  Score: %f\n", *entity.Score)
				fmt.Printf("  Type:  %v\n", entity.Type)
				fmt.Printf("  Text:  %s\n", *entity.Text)
				fmt.Printf("  BegOffset: %d\n", *entity.BeginOffset)
				fmt.Printf("  EndOffest: %d\n", *entity.EndOffset)
				fmt.Printf("  Types:     %v\n", entity.Type.Values())
				fmt.Println("---------------------------------------")
			}

		}
	}
	return nil
}

func runSentiment(text []string) error {
	resp, err := c.Sentiment(&text)
	if err != nil {
		log.Fatal(err)
	}
	if len(resp.ErrorList) > 0 {
		for _, err := range resp.ErrorList {
			log.Error(err)
		}
	}
	for _, item := range resp.ResultList {
		fmt.Printf("%d :: Text: %s\n", *item.Index, text[*item.Index])
		fmt.Printf("  Sentiment: %v\n", item.Sentiment)
		fmt.Printf("  Mixed: %f\n", *item.SentimentScore.Mixed)
		fmt.Printf("  Negative: %f\n", *item.SentimentScore.Negative)
		fmt.Printf("  Neutral: %f\n", *item.SentimentScore.Neutral)
		fmt.Printf("  Positive: %f\n", *item.SentimentScore.Positive)
		fmt.Println("----------------------")
	}
	return nil
}
