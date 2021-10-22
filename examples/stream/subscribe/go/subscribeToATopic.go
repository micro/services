package example

import (
	"fmt"
	"github.com/micro/services/clients/go/stream"
	"os"
)

// Subscribe to messages for a given topic.
func SubscribeToAtopic() {
	streamService := stream.NewStreamService(os.Getenv("MICRO_API_TOKEN"))
	rsp, err := streamService.Subscribe(&stream.SubscribeRequest{
		Topic: "events",
	})
	fmt.Println(rsp, err)
}
