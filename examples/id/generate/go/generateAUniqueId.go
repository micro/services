package example

import (
	"fmt"
	"github.com/micro/services/clients/go/id"
	"os"
)

// Generate a unique ID. Defaults to uuid.
func GenerateAuniqueId() {
	idService := id.NewIdService(os.Getenv("MICRO_API_TOKEN"))
	rsp, err := idService.Generate(&id.GenerateRequest{
		Type: "uuid",
	})
	fmt.Println(rsp, err)
}
