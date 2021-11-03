package example

import (
	"fmt"
	"os"

	"github.com/micro/services/clients/go/stream"
)

// List messages for a given channel
func ListMessages() {
	streamService := stream.NewStreamService(os.Getenv("MICRO_API_TOKEN"))
	rsp, err := streamService.ListMessages(&stream.ListMessagesRequest{
		Channel: "general",
	})
	fmt.Println(rsp, err)
}
