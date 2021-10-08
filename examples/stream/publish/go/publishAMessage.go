package example

import (
	"fmt"
	"github.com/micro/services/clients/go/stream"
	"os"
)

// Publish a message to the stream. Specify a topic to group messages for a specific topic.
func PublishAmessage() {
	streamService := stream.NewStreamService(os.Getenv("MICRO_API_TOKEN"))
	rsp, err := streamService.Publish(&stream.PublishRequest{
		Message: map[string]interface{}{
			"id":   "1",
			"type": "signup",
			"user": "john",
		},
		Topic: "events",
	})
	fmt.Println(rsp, err)
}
