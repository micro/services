package example

import (
	"fmt"
	"os"

	"github.com/micro/services/clients/go/function"
)

// Call a function by name
func CallAfunction() {
	functionService := function.NewFunctionService(os.Getenv("MICRO_API_TOKEN"))
	rsp, err := functionService.Call(&function.CallRequest{
		Name:    "my-first-func",
		Request: map[string]interface{}{},
	})
	fmt.Println(rsp, err)
}
