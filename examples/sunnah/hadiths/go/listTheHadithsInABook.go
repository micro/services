package example

import (
	"fmt"
	"os"

	"github.com/micro/services/clients/go/sunnah"
)

// Hadiths returns a list of hadiths and their corresponding text for a
// given book within a collection.
func ListTheHadithsInAbook() {
	sunnahService := sunnah.NewSunnahService(os.Getenv("MICRO_API_TOKEN"))
	rsp, err := sunnahService.Hadiths(&sunnah.HadithsRequest{
		Book:       1,
		Collection: "bukhari",
	})
	fmt.Println(rsp, err)
}
