package example

import (
	"fmt"
	"os"

	"github.com/micro/services/clients/go/id"
)

// Generate a unique ID. Defaults to uuid.
func GenerateAsnowflakeId() {
	idService := id.NewIdService(os.Getenv("MICRO_API_TOKEN"))
	rsp, err := idService.Generate(&id.GenerateRequest{
		Type: "snowflake",
	})
	fmt.Println(rsp, err)
}
