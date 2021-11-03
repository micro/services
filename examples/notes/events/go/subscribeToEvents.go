package example

import (
	"fmt"
	"os"

	"github.com/micro/services/clients/go/notes"
)

// Subscribe to notes events
func SubscribeToEvents() {
	notesService := notes.NewNotesService(os.Getenv("MICRO_API_TOKEN"))
	rsp, err := notesService.Events(&notes.EventsRequest{
		Id: "63c0cdf8-2121-11ec-a881-0242e36f037a",
	})
	fmt.Println(rsp, err)
}
