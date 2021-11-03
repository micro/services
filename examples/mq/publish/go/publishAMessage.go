package example

import (
	"fmt"
	"os"

	"github.com/micro/services/clients/go/mq"
)

// Publish a message. Specify a topic to group messages for a specific topic.
func PublishAmessage() {
	mqService := mq.NewMqService(os.Getenv("MICRO_API_TOKEN"))
	rsp, err := mqService.Publish(&mq.PublishRequest{
		Message: map[string]interface{}{
			"user": "john",
			"id":   "1",
			"type": "signup",
		},
		Topic: "events",
	})
	fmt.Println(rsp, err)
}
