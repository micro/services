package example

import (
	"fmt"
	"os"

	"github.com/micro/services/clients/go/stream"
)

// SendMessage a message to the stream.
func SendAmessage() {
	streamService := stream.NewStreamService(os.Getenv("MICRO_API_TOKEN"))
	rsp, err := streamService.SendMessage(&stream.SendMessageRequest{
		Channel: "general",
		Text:    "Hey checkout this tweet https://twitter.com/m3oservices/status/1455291054295498752",
	})
	fmt.Println(rsp, err)
}
