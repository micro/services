package example

import (
	"fmt"
	"os"

	"github.com/micro/services/clients/go/url"
)

// Shortens a destination URL and returns a full short URL.
func ShortenAlongUrl() {
	urlService := url.NewUrlService(os.Getenv("MICRO_API_TOKEN"))
	rsp, err := urlService.Shorten(&url.ShortenRequest{})
	fmt.Println(rsp, err)
}
