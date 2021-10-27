# Time

An [m3o.com](https://m3o.com) API. For example usage see [m3o.com/Time/api](https://m3o.com/Time/api).

Endpoints:

## Now

Get the current time


[https://m3o.com/time/api#Now](https://m3o.com/time/api#Now)

```go
package example

import(
	"fmt"
	"os"

	"github.com/micro/services/clients/go/time"
)

// Get the current time
func ReturnsCurrentTimeOptionallyWithLocation() {
	timeService := time.NewTimeService(os.Getenv("MICRO_API_TOKEN"))
	rsp, err := timeService.Now(&time.NowRequest{
		
	})
	fmt.Println(rsp, err)
}
```
## Zone

Get the timezone info for a specific location


[https://m3o.com/time/api#Zone](https://m3o.com/time/api#Zone)

```go
package example

import(
	"fmt"
	"os"

	"github.com/micro/services/clients/go/time"
)

// Get the timezone info for a specific location
func GetTheTimezoneInfoForAspecificLocation() {
	timeService := time.NewTimeService(os.Getenv("MICRO_API_TOKEN"))
	rsp, err := timeService.Zone(&time.ZoneRequest{
		Location: "London",

	})
	fmt.Println(rsp, err)
}
```
