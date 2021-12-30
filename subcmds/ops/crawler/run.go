package crawler

import "fmt"

func runRunCrawler() error {
	if err := crawler.StartCrawler(); err != nil {
		return err
	}
	fmt.Println("Crawler started")
	return nil
}
