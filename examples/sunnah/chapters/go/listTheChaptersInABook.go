package example

import (
	"fmt"
	"os"

	"github.com/micro/services/clients/go/sunnah"
)

// Get all the chapters of a given book within a collection.
func ListTheChaptersInAbook() {
	sunnahService := sunnah.NewSunnahService(os.Getenv("MICRO_API_TOKEN"))
	rsp, err := sunnahService.Chapters(&sunnah.ChaptersRequest{
		Book:       1,
		Collection: "bukhari",
	})
	fmt.Println(rsp, err)
}
