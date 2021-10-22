package example

import (
	"fmt"
	"os"

	"github.com/micro/services/clients/go/rss"
)

// Add a new RSS feed with a name, url, and category
func AddAnewFeed() {
	rssService := rss.NewRssService(os.Getenv("MICRO_API_TOKEN"))
	rsp, err := rssService.Add(&rss.AddRequest{
		Category: "news",
		Name:     "bbc",
		Url:      "http://feeds.bbci.co.uk/news/rss.xml",
	})
	fmt.Println(rsp, err)
}
