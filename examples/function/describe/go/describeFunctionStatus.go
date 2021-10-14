package example

import (
	"fmt"
	"github.com/micro/services/clients/go/function"
	"os"
)

//
func DescribeFunctionStatus() {
	functionService := function.NewFunctionService(os.Getenv("MICRO_API_TOKEN"))
	rsp, err := functionService.Describe(&function.DescribeRequest{
		Name:    "my-first-func",
		Project: "tests",
	})
	fmt.Println(rsp, err)
}
