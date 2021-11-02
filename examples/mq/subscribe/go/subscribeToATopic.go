package example

import (
	"fmt"
	"os"

	"github.com/micro/services/clients/go/mq"
)

// Subscribe to messages for a given topic.
func SubscribeToAtopic() {
	mqService := mq.NewMqService(os.Getenv("MICRO_API_TOKEN"))
	rsp, err := mqService.Subscribe(&mq.SubscribeRequest{
		Topic: "events",
	})
	fmt.Println(rsp, err)
}
