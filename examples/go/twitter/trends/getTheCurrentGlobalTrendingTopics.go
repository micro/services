package example

import (
	"fmt"
	"os"

	"github.com/micro/services/clients/go/twitter"
)

// Get the current global trending topics
func GetTheCurrentGlobalTrendingTopics() {
	twitterService := twitter.NewTwitterService(os.Getenv("MICRO_API_TOKEN"))
	rsp, err := twitterService.Trends(&twitter.TrendsRequest{})
	fmt.Println(rsp, err)
}
