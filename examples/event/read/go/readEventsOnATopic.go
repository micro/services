package example

import (
	"fmt"
	"os"

	"github.com/micro/services/clients/go/event"
)

// Read stored events
func ReadEventsOnAtopic() {
	eventService := event.NewEventService(os.Getenv("MICRO_API_TOKEN"))
	rsp, err := eventService.Read(&event.ReadRequest{
		Topic: "user",
	})
	fmt.Println(rsp, err)
}
