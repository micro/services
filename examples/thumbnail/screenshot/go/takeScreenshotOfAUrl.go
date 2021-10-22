package example

import (
	"fmt"
	"os"

	"github.com/micro/services/clients/go/thumbnail"
)

// Create a thumbnail screenshot by passing in a url, height and width
func TakeScreenshotOfAurl() {
	thumbnailService := thumbnail.NewThumbnailService(os.Getenv("MICRO_API_TOKEN"))
	rsp, err := thumbnailService.Screenshot(&thumbnail.ScreenshotRequest{
		Height: 600,
		Url:    "https://m3o.com",
		Width:  600,
	})
	fmt.Println(rsp, err)
}
