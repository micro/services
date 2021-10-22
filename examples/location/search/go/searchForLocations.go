package example

import (
	"fmt"
	"os"

	"github.com/micro/services/clients/go/location"
)

// Search for entities in a given radius
func SearchForLocations() {
	locationService := location.NewLocationService(os.Getenv("MICRO_API_TOKEN"))
	rsp, err := locationService.Search(&location.SearchRequest{
		Center: &location.Point{
			Latitude:  51.511061,
			Longitude: -0.120022,
		},
		NumEntities: 10,
		Radius:      100,
		Type:        "bike",
	})
	fmt.Println(rsp, err)
}
