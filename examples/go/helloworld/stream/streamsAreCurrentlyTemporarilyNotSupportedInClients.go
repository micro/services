package example

import (
	"fmt"
	"os"

	"github.com/micro/services/clients/go/helloworld"
)

// Stream returns a stream of "Hello $name" responses
func StreamsAreCurrentlyTemporarilyNotSupportedInClients() {
	helloworldService := helloworld.NewHelloworldService(os.Getenv("MICRO_API_TOKEN"))
	rsp, err := helloworldService.Stream(&helloworld.StreamRequest{
		Name: "not supported",
	})
	fmt.Println(rsp, err)
}
