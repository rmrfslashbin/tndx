package crawler

import (
	"fmt"
)

func runCrawlerStatus() error {
	if ret, err := crawler.GetCrawlerData(); err != nil {
		return err
	} else {
		fmt.Println("Crawler Status")
		fmt.Printf("Name:                   %s\n", *ret.Crawler.Name)
		fmt.Printf("State:                  %s\n", ret.Crawler.State)
		fmt.Printf("Elapsed Time:           %d (seconds)\n", ret.Crawler.CrawlElapsedTime/1000)
		fmt.Printf("Last Crawl Status:      %v\n", ret.Crawler.LastCrawl.Status)
		fmt.Printf("Last Crawl Error:       %v\n", ret.Crawler.LastCrawl.ErrorMessage)
		fmt.Printf("Last Crawl Start time:  %v\n", ret.Crawler.LastCrawl.StartTime)
		fmt.Println()
	}
	return nil
}
