package example

import (
	"fmt"
	"github.com/micro/services/clients/go/function"
	"os"
)

// Deploy a group of functions
func DeployAfunction() {
	functionService := function.NewFunctionService(os.Getenv("MICRO_API_TOKEN"))
	rsp, err := functionService.Deploy(&function.DeployRequest{
		Entrypoint: "helloworld",
		Name:       "my-first-func",
		Project:    "tests",
		Repo:       "github.com/m3o/nodejs-function-example",
		Runtime:    "nodejs14",
	})
	fmt.Println(rsp, err)
}
