package example

import (
	"fmt"
	"os"

	"github.com/micro/services/clients/go/stream"
)

// List all the active channels
func ListChannels() {
	streamService := stream.NewStreamService(os.Getenv("MICRO_API_TOKEN"))
	rsp, err := streamService.ListChannels(&stream.ListChannelsRequest{})
	fmt.Println(rsp, err)
}
