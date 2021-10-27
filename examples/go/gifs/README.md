# Gifs

An [m3o.com](https://m3o.com) API. For example usage see [m3o.com/Gifs/api](https://m3o.com/Gifs/api).

Endpoints:

## Search

Search for a GIF


[https://m3o.com/gifs/api#Search](https://m3o.com/gifs/api#Search)

```go
package example

import(
	"fmt"
	"os"

	"github.com/micro/services/clients/go/gifs"
)

// Search for a GIF
func Search() {
	gifsService := gifs.NewGifsService(os.Getenv("MICRO_API_TOKEN"))
	rsp, err := gifsService.Search(&gifs.SearchRequest{
		Limit: 2,
Query: "dogs",

	})
	fmt.Println(rsp, err)
}
```
