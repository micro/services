package example

import (
	"fmt"
	"os"

	"github.com/micro/services/clients/go/rss"
)

// Get an RSS feed by name. If no name is given, all feeds are returned. Default limit is 25 entries.
func ReadAfeed() {
	rssService := rss.NewRssService(os.Getenv("MICRO_API_TOKEN"))
	rsp, err := rssService.Feed(&rss.FeedRequest{
		Name: "bbc",
	})
	fmt.Println(rsp, err)
}
