package example

import (
	"fmt"
	"github.com/micro/services/clients/go/postcode"
	"os"
)

// Lookup a postcode to retrieve the related region, county, etc
func LookupPostcode() {
	postcodeService := postcode.NewPostcodeService(os.Getenv("MICRO_API_TOKEN"))
	rsp, err := postcodeService.Lookup(&postcode.LookupRequest{
		Postcode: "SW1A 2AA",
	})
	fmt.Println(rsp, err)
}
