package example

import (
	"fmt"
	"os"

	"github.com/micro/services/clients/go/google"
)

// Search for videos on Google
func SearchForVideos() {
	googleService := google.NewGoogleService(os.Getenv("MICRO_API_TOKEN"))
	rsp, err := googleService.Search(&google.SearchRequest{
		Query: "how to make donuts",
	})
	fmt.Println(rsp, err)
}
