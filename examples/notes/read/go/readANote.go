package example

import (
	"fmt"
	"github.com/micro/services/clients/go/notes"
	"os"
)

// Read a note
func ReadAnote() {
	notesService := notes.NewNotesService(os.Getenv("MICRO_API_TOKEN"))
	rsp, err := notesService.Read(&notes.ReadRequest{
		Id: "63c0cdf8-2121-11ec-a881-0242e36f037a",
	})
	fmt.Println(rsp, err)
}
