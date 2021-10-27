# Url

An [m3o.com](https://m3o.com) API. For example usage see [m3o.com/Url/api](https://m3o.com/Url/api).

Endpoints:

## List

List information on all the shortened URLs that you have created


[https://m3o.com/url/api#List](https://m3o.com/url/api#List)

```go
package example

import(
	"fmt"
	"os"

	"github.com/micro/services/clients/go/url"
)

// List information on all the shortened URLs that you have created
func ListYourShortenedUrls() {
	urlService := url.NewUrlService(os.Getenv("MICRO_API_TOKEN"))
	rsp, err := urlService.List(&url.ListRequest{
		
	})
	fmt.Println(rsp, err)
}
```
## Shorten

Shortens a destination URL and returns a full short URL.


[https://m3o.com/url/api#Shorten](https://m3o.com/url/api#Shorten)

```go
package example

import(
	"fmt"
	"os"

	"github.com/micro/services/clients/go/url"
)

// Shortens a destination URL and returns a full short URL.
func ShortenAlongUrl() {
	urlService := url.NewUrlService(os.Getenv("MICRO_API_TOKEN"))
	rsp, err := urlService.Shorten(&url.ShortenRequest{
		
	})
	fmt.Println(rsp, err)
}
```
## Proxy

Proxy returns the destination URL of a short URL.


[https://m3o.com/url/api#Proxy](https://m3o.com/url/api#Proxy)

```go
package example

import(
	"fmt"
	"os"

	"github.com/micro/services/clients/go/url"
)

// Proxy returns the destination URL of a short URL.
func ResolveAshortUrlToAlongDestinationUrl() {
	urlService := url.NewUrlService(os.Getenv("MICRO_API_TOKEN"))
	rsp, err := urlService.Proxy(&url.ProxyRequest{
		
	})
	fmt.Println(rsp, err)
}
```
