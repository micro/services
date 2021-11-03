package example

import (
	"fmt"
	"os"

	"github.com/micro/services/clients/go/event"
)

// Subscribe to messages for a given topic.
func SubscribeToAtopic() {
	eventService := event.NewEventService(os.Getenv("MICRO_API_TOKEN"))
	rsp, err := eventService.Subscribe(&event.SubscribeRequest{
		Topic: "user",
	})
	fmt.Println(rsp, err)
}
