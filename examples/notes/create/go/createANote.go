package example

import (
	"fmt"
	"github.com/micro/services/clients/go/notes"
	"os"
)

// Create a new note
func CreateAnote() {
	notesService := notes.NewNotesService(os.Getenv("MICRO_API_TOKEN"))
	rsp, err := notesService.Create(&notes.CreateRequest{
		Text:  "This is my note",
		Title: "New Note",
	})
	fmt.Println(rsp, err)
}
