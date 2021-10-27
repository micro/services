package example

import (
	"fmt"
	"os"

	"github.com/micro/services/clients/go/routing"
)

// Get the eta for a route from origin to destination. The eta is an estimated time based on car routes
func EtaFromPointAtoPointB() {
	routingService := routing.NewRoutingService(os.Getenv("MICRO_API_TOKEN"))
	rsp, err := routingService.Eta(&routing.EtaRequest{
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
