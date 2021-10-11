package example

import (
	"fmt"
	"github.com/micro/services/clients/go/function"
	"os"
)

//
func ListFunctions() {
	functionService := function.NewFunctionService(os.Getenv("MICRO_API_TOKEN"))
	rsp, err := functionService.List(&function.ListRequest{})
	fmt.Println(rsp, err)
}
