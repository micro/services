package example

import (
	"fmt"
	"os"

	"github.com/micro/services/clients/go/helloworld"
)

// Stream returns a stream of "Hello $name" responses
func StreamsResponsesFromTheServerUsingWebsockets() {
	helloworldService := helloworld.NewHelloworldService(os.Getenv("MICRO_API_TOKEN"))
	rsp, err := helloworldService.Stream(&helloworld.StreamRequest{
		Messages: 10,
		Name:     "John",
	})
	fmt.Println(rsp, err)
}
