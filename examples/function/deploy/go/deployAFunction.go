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
		Repo:       "github.com/crufter/gcloud-nodejs-test",
	})
	fmt.Println(rsp, err)
}
