package example

import (
	"fmt"
	"os"

	"github.com/micro/services/clients/go/event"
)

// Consume events from a given topic.
func ConsumeFromAtopic() {
	eventService := event.NewEventService(os.Getenv("MICRO_API_TOKEN"))
	rsp, err := eventService.Consume(&event.ConsumeRequest{
		Topic: "user",
	})
	fmt.Println(rsp, err)
}
