# Twitter

An [m3o.com](https://m3o.com) API. For example usage see [m3o.com/Twitter/api](https://m3o.com/Twitter/api).

Endpoints:

## Search

Search for tweets with a simple query


[https://m3o.com/twitter/api#Search](https://m3o.com/twitter/api#Search)

```go
package example

import(
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
```
## Trends

Get the current global trending topics


[https://m3o.com/twitter/api#Trends](https://m3o.com/twitter/api#Trends)

```go
package example

import(
	"fmt"
	"os"

	"github.com/micro/services/clients/go/twitter"
)

// Get the current global trending topics
func GetTheCurrentGlobalTrendingTopics() {
	twitterService := twitter.NewTwitterService(os.Getenv("MICRO_API_TOKEN"))
	rsp, err := twitterService.Trends(&twitter.TrendsRequest{
		
	})
	fmt.Println(rsp, err)
}
```
## User

Get a user's twitter profile


[https://m3o.com/twitter/api#User](https://m3o.com/twitter/api#User)

```go
package example

import(
	"fmt"
	"os"

	"github.com/micro/services/clients/go/twitter"
)

// Get a user's twitter profile
func GetAusersTwitterProfile() {
	twitterService := twitter.NewTwitterService(os.Getenv("MICRO_API_TOKEN"))
	rsp, err := twitterService.User(&twitter.UserRequest{
		Username: "crufter",

	})
	fmt.Println(rsp, err)
}
```
## Timeline

Get the timeline for a given user


[https://m3o.com/twitter/api#Timeline](https://m3o.com/twitter/api#Timeline)

```go
package example

import(
	"fmt"
	"os"

	"github.com/micro/services/clients/go/twitter"
)

// Get the timeline for a given user
func GetAtwitterTimeline() {
	twitterService := twitter.NewTwitterService(os.Getenv("MICRO_API_TOKEN"))
	rsp, err := twitterService.Timeline(&twitter.TimelineRequest{
		Limit: 1,
Username: "m3oservices",

	})
	fmt.Println(rsp, err)
}
```
