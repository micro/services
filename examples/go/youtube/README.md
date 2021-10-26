# Youtube

An [m3o.com](https://m3o.com) API. For example usage see [m3o.com/Youtube/api](https://m3o.com/Youtube/api).

Endpoints:

## Search

Search for videos on YouTube


[https://m3o.com/youtube/api#Search](https://m3o.com/youtube/api#Search)

```go
package example

import(
	"fmt"
	"os"

	"github.com/micro/services/clients/go/youtube"
)

// Search for videos on YouTube
func SearchForVideos() {
	youtubeService := youtube.NewYoutubeService(os.Getenv("MICRO_API_TOKEN"))
	rsp, err := youtubeService.Search(&youtube.SearchRequest{
		Query: "donuts",

	})
	fmt.Println(rsp, err)
}
```
