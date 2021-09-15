package example

import (
	"fmt"
	"github.com/micro/services/clients/go/id"
	"os"
)

// Generate a unique ID. Defaults to uuid.
func GenerateAbigflakeId() {
	idService := id.NewIdService(os.Getenv("MICRO_API_TOKEN"))
	rsp, err := idService.Generate(&id.GenerateRequest{
		Type: "bigflake",
	})
	fmt.Println(rsp, err)
}
