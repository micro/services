package example

import (
	"fmt"
	"os"

	"github.com/micro/services/clients/go/twitter"
)

// Search for tweets with a simple query
func SearchForTweets() {
	twitterService := twitter.NewTwitterService(os.Getenv("MICRO_API_TOKEN"))
	rsp, err := twitterService.Search(&twitter.SearchRequest{
		Query: "cats",
	})
	fmt.Println(rsp, err)
}
