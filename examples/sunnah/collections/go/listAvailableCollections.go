package example

import (
	"fmt"
	"github.com/micro/services/clients/go/sunnah"
	"os"
)

// Get a list of available collections. A collection is
// a compilation of hadiths collected and written by an author.
func ListAvailableCollections() {
	sunnahService := sunnah.NewSunnahService(os.Getenv("MICRO_API_TOKEN"))
	rsp, err := sunnahService.Collections(&sunnah.CollectionsRequest{})
	fmt.Println(rsp, err)
}
