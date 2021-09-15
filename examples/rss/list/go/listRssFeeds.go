package example

import (
	"fmt"
	"github.com/micro/services/clients/go/rss"
	"os"
)

// List the saved RSS fields
func ListRssFeeds() {
	rssService := rss.NewRssService(os.Getenv("MICRO_API_TOKEN"))
	rsp, err := rssService.List(&rss.ListRequest{})
	fmt.Println(rsp, err)
}
