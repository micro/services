package example

import (
	"fmt"
	"os"

	"github.com/micro/services/clients/go/weather"
)

// Get the weather forecast for the next 1-10 days
func ForecastWeather() {
	weatherService := weather.NewWeatherService(os.Getenv("MICRO_API_TOKEN"))
	rsp, err := weatherService.Forecast(&weather.ForecastRequest{
		Days:     2,
		Location: "London",
	})
	fmt.Println(rsp, err)
}
