# Ip

An [m3o.com](https://m3o.com) API. For example usage see [m3o.com/Ip/api](https://m3o.com/Ip/api).

Endpoints:

## Lookup

Lookup the geolocation information for an IP address


[https://m3o.com/ip/api#Lookup](https://m3o.com/ip/api#Lookup)

```go
package example

import(
	"fmt"
	"os"

	"github.com/micro/services/clients/go/ip"
)

// Lookup the geolocation information for an IP address
func LookupIpInfo() {
	ipService := ip.NewIpService(os.Getenv("MICRO_API_TOKEN"))
	rsp, err := ipService.Lookup(&ip.LookupRequest{
		Ip: "93.148.214.31",

	})
	fmt.Println(rsp, err)
}
```
