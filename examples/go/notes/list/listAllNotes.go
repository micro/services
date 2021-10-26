package example

import (
	"fmt"
	"os"

	"github.com/micro/services/clients/go/notes"
)

// List all the notes
func ListAllNotes() {
	notesService := notes.NewNotesService(os.Getenv("MICRO_API_TOKEN"))
	rsp, err := notesService.List(&notes.ListRequest{})
	fmt.Println(rsp, err)
}
