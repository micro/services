package example

import (
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
