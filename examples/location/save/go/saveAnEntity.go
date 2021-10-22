package example

import (
	"fmt"
	"os"

	"github.com/micro/services/clients/go/location"
)

// Save an entity's current position
func SaveAnEntity() {
	locationService := location.NewLocationService(os.Getenv("MICRO_API_TOKEN"))
	rsp, err := locationService.Save(&location.SaveRequest{
		Entity: &location.Entity{
			Id: "1",
			Location: &location.Point{
				Latitude:  51.511061,
				Longitude: -0.120022,
				Timestamp: 1622802761,
			},
			Type: "bike",
		},
	})
	fmt.Println(rsp, err)
}
