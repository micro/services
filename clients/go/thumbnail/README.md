# Thumbnail

An [m3o.com](https://m3o.com) API. For example usage see [m3o.com/Thumbnail/api](https://m3o.com/Thumbnail/api).

Endpoints:

## Screenshot

Create a thumbnail screenshot by passing in a url, height and width


[https://m3o.com/thumbnail/api#Screenshot](https://m3o.com/thumbnail/api#Screenshot)

```go
package example

import(
	"fmt"
	"os"

	"github.com/micro/services/clients/go/thumbnail"
)

// Create a thumbnail screenshot by passing in a url, height and width
func TakeScreenshotOfAurl() {
	thumbnailService := thumbnail.NewThumbnailService(os.Getenv("MICRO_API_TOKEN"))
	rsp, err := thumbnailService.Screenshot(&thumbnail.ScreenshotRequest{
		Height: 600,
Url: "https://m3o.com",
Width: 600,

	})
	fmt.Println(rsp, err)
}
```
