package example

import (
	"fmt"
	"github.com/micro/services/clients/go/rss"
	"os"
)

// Remove an RSS feed by name
func RemoveAfeed() {
	rssService := rss.NewRssService(os.Getenv("MICRO_API_TOKEN"))
	rsp, err := rssService.Remove(&rss.RemoveRequest{
		Name: "bbc",
	})
	fmt.Println(rsp, err)
}
