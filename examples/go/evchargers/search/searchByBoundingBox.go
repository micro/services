package example

import (
	"fmt"
	"os"

	"github.com/micro/services/clients/go/evchargers"
)

// Search by giving a coordinate and a max distance, or bounding box and optional filters
func SearchByBoundingBox() {
	evchargersService := evchargers.NewEvchargersService(os.Getenv("MICRO_API_TOKEN"))
	rsp, err := evchargersService.Search(&evchargers.SearchRequest{
		Box: &evchargers.BoundingBox{},
	})
	fmt.Println(rsp, err)
}
