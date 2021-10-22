package example

import (
	"fmt"
	"os"

	"github.com/micro/services/clients/go/id"
)

// List the types of IDs available. No query params needed.
func ListTheTypesOfIdsAvailable() {
	idService := id.NewIdService(os.Getenv("MICRO_API_TOKEN"))
	rsp, err := idService.Types(&id.TypesRequest{})
	fmt.Println(rsp, err)
}
