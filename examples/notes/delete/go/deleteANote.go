package example

import (
	"fmt"
	"os"

	"github.com/micro/services/clients/go/notes"
)

// Delete a note
func DeleteAnote() {
	notesService := notes.NewNotesService(os.Getenv("MICRO_API_TOKEN"))
	rsp, err := notesService.Delete(&notes.DeleteRequest{
		Id: "63c0cdf8-2121-11ec-a881-0242e36f037a",
	})
	fmt.Println(rsp, err)
}
