# Prayer

An [m3o.com](https://m3o.com) API. For example usage see [m3o.com/Prayer/api](https://m3o.com/Prayer/api).

Endpoints:

## Times

Get the prayer (salah) times for a location on a given date


[https://m3o.com/prayer/api#Times](https://m3o.com/prayer/api#Times)

```go
package example

import(
	"fmt"
	"os"

	"github.com/micro/services/clients/go/prayer"
)

// Get the prayer (salah) times for a location on a given date
func PrayerTimes() {
	prayerService := prayer.NewPrayerService(os.Getenv("MICRO_API_TOKEN"))
	rsp, err := prayerService.Times(&prayer.TimesRequest{
		Location: "london",

	})
	fmt.Println(rsp, err)
}
```
