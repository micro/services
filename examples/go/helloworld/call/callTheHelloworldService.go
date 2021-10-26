package example

import (
	"fmt"
	"os"

	"github.com/micro/services/clients/go/helloworld"
)

// Call returns a personalised "Hello $name" response
func CallTheHelloworldService() {
	helloworldService := helloworld.NewHelloworldService(os.Getenv("MICRO_API_TOKEN"))
	rsp, err := helloworldService.Call(&helloworld.CallRequest{
		Name: "John",
	})
	fmt.Println(rsp, err)
}
