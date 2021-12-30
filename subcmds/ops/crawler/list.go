package crawler

import (
	"fmt"
)

func runCrawlerList() error {
	if ret, err := crawler.ListCrawlers(); err != nil {
		return err
	} else {
		fmt.Println("Crawler List:")
		for _, name := range ret {
			fmt.Printf("  %s\n", name)
		}
		fmt.Println()
	}
	return nil
}
