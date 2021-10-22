package example

import (
	"fmt"
	"os"

	"github.com/micro/services/clients/go/function"
)

// List all the deployed functions
func ListFunctions() {
	functionService := function.NewFunctionService(os.Getenv("MICRO_API_TOKEN"))
	rsp, err := functionService.List(&function.ListRequest{})
	fmt.Println(rsp, err)
}
