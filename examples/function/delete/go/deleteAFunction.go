package example

import (
	"fmt"
	"github.com/micro/services/clients/go/function"
	"os"
)

// Delete a function by name
func DeleteAfunction() {
	functionService := function.NewFunctionService(os.Getenv("MICRO_API_TOKEN"))
	rsp, err := functionService.Delete(&function.DeleteRequest{
		Name:    "my-first-func",
		Project: "tests",
	})
	fmt.Println(rsp, err)
}
