package example

import (
	"fmt"
	"os"

	"github.com/micro/services/clients/go/url"
)

// Proxy returns the destination URL of a short URL.
func ResolveAshortUrlToAlongDestinationUrl() {
	urlService := url.NewUrlService(os.Getenv("MICRO_API_TOKEN"))
	rsp, err := urlService.Proxy(&url.ProxyRequest{})
	fmt.Println(rsp, err)
}
