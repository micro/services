package example

import (
	"fmt"
	"os"

	"github.com/micro/services/clients/go/url"
)

// List information on all the shortened URLs that you have created
func ListYourShortenedUrls() {
	urlService := url.NewUrlService(os.Getenv("MICRO_API_TOKEN"))
	rsp, err := urlService.List(&url.ListRequest{})
	fmt.Println(rsp, err)
}
