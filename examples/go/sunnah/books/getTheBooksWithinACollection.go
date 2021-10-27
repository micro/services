package example

import (
	"fmt"
	"os"

	"github.com/micro/services/clients/go/sunnah"
)

// Get a list of books from within a collection. A book can contain many chapters
// each with its own hadiths.
func GetTheBooksWithinAcollection() {
	sunnahService := sunnah.NewSunnahService(os.Getenv("MICRO_API_TOKEN"))
	rsp, err := sunnahService.Books(&sunnah.BooksRequest{
		Collection: "bukhari",
	})
	fmt.Println(rsp, err)
}
