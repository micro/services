package example

import (
	"fmt"
	"os"

	"github.com/micro/services/clients/go/event"
)

// Publish a event to the event stream.
func PublishAnEvent() {
	eventService := event.NewEventService(os.Getenv("MICRO_API_TOKEN"))
	rsp, err := eventService.Publish(&event.PublishRequest{
		Message: map[string]interface{}{
			"id":   "1",
			"type": "signup",
			"user": "john",
		},
		Topic: "user",
	})
	fmt.Println(rsp, err)
}
