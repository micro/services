package example

import (
	"fmt"
	"github.com/micro/services/clients/go/url"
	"os"
)

// Proxy returns the destination URL of a short URL.
func ResolveAshortUrlToAlongDestinationUrl() {
	urlService := url.NewUrlService(os.Getenv("MICRO_API_TOKEN"))
	rsp, err := urlService.Proxy(&url.ProxyRequest{})
	fmt.Println(rsp, err)
}
