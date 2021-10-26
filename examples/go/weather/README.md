# Weather

An [m3o.com](https://m3o.com) API. For example usage see [m3o.com/Weather/api](https://m3o.com/Weather/api).

Endpoints:

## Now

Get the current weather report for a location by postcode, city, zip code, ip address


[https://m3o.com/weather/api#Now](https://m3o.com/weather/api#Now)

```go
package example

import(
	"fmt"
	"os"

	"github.com/micro/services/clients/go/weather"
)

// Get the current weather report for a location by postcode, city, zip code, ip address
func GetCurrentWeather() {
	weatherService := weather.NewWeatherService(os.Getenv("MICRO_API_TOKEN"))
	rsp, err := weatherService.Now(&weather.NowRequest{
		Location: "london",

	})
	fmt.Println(rsp, err)
}
```
## Forecast

Get the weather forecast for the next 1-10 days


[https://m3o.com/weather/api#Forecast](https://m3o.com/weather/api#Forecast)

```go
package example

import(
	"fmt"
	"os"

	"github.com/micro/services/clients/go/weather"
)

// Get the weather forecast for the next 1-10 days
func ForecastWeather() {
	weatherService := weather.NewWeatherService(os.Getenv("MICRO_API_TOKEN"))
	rsp, err := weatherService.Forecast(&weather.ForecastRequest{
		Days: 2,
Location: "London",

	})
	fmt.Println(rsp, err)
}
```
