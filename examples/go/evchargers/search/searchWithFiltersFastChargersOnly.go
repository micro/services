package example

import (
	"fmt"
	"os"

	"github.com/micro/services/clients/go/evchargers"
)

// Search by giving a coordinate and a max distance, or bounding box and optional filters
func SearchWithFiltersFastChargersOnly() {
	evchargersService := evchargers.NewEvchargersService(os.Getenv("MICRO_API_TOKEN"))
	rsp, err := evchargersService.Search(&evchargers.SearchRequest{
		Distance: 2000,
		Levels:   []string{"3"},
		Location: &evchargers.Coordinates{
			Latitude:  51.53336351319885,
			Longitude: -0.0252,
		},
	})
	fmt.Println(rsp, err)
}
