package example

import (
	"fmt"
	"os"

	"github.com/micro/services/clients/go/routing"
)

// Retrieve a route as a simple list of gps points along with total distance and estimated duration
func GpsPointsForAroute() {
	routingService := routing.NewRoutingService(os.Getenv("MICRO_API_TOKEN"))
	rsp, err := routingService.Route(&routing.RouteRequest{
		Destination: &routing.Point{
			Latitude:  52.529407,
			Longitude: 13.397634,
		},
		Origin: &routing.Point{
			Latitude:  52.517037,
			Longitude: 13.38886,
		},
	})
	fmt.Println(rsp, err)
}
