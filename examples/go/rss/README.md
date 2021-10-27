# Rss

An [m3o.com](https://m3o.com) API. For example usage see [m3o.com/Rss/api](https://m3o.com/Rss/api).

Endpoints:

## List

List the saved RSS fields


[https://m3o.com/rss/api#List](https://m3o.com/rss/api#List)

```go
package example

import(
	"fmt"
	"os"

	"github.com/micro/services/clients/go/rss"
)

// List the saved RSS fields
func ListRssFeeds() {
	rssService := rss.NewRssService(os.Getenv("MICRO_API_TOKEN"))
	rsp, err := rssService.List(&rss.ListRequest{
		
	})
	fmt.Println(rsp, err)
}
```
## Remove

Remove an RSS feed by name


[https://m3o.com/rss/api#Remove](https://m3o.com/rss/api#Remove)

```go
package example

import(
	"fmt"
	"os"

	"github.com/micro/services/clients/go/rss"
)

// Remove an RSS feed by name
func RemoveAfeed() {
	rssService := rss.NewRssService(os.Getenv("MICRO_API_TOKEN"))
	rsp, err := rssService.Remove(&rss.RemoveRequest{
		Name: "bbc",

	})
	fmt.Println(rsp, err)
}
```
## Add

Add a new RSS feed with a name, url, and category


[https://m3o.com/rss/api#Add](https://m3o.com/rss/api#Add)

```go
package example

import(
	"fmt"
	"os"

	"github.com/micro/services/clients/go/rss"
)

// Add a new RSS feed with a name, url, and category
func AddAnewFeed() {
	rssService := rss.NewRssService(os.Getenv("MICRO_API_TOKEN"))
	rsp, err := rssService.Add(&rss.AddRequest{
		Category: "news",
Name: "bbc",
Url: "http://feeds.bbci.co.uk/news/rss.xml",

	})
	fmt.Println(rsp, err)
}
```
## Feed

Get an RSS feed by name. If no name is given, all feeds are returned. Default limit is 25 entries.


[https://m3o.com/rss/api#Feed](https://m3o.com/rss/api#Feed)

```go
package example

import(
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
```
