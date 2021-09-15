package example

import (
	"fmt"
	"github.com/micro/services/clients/go/postcode"
	"os"
)

// Validate a postcode.
func ReturnArandomPostcodeAndItsInformation() {
	postcodeService := postcode.NewPostcodeService(os.Getenv("MICRO_API_TOKEN"))
	rsp, err := postcodeService.Validate(&postcode.ValidateRequest{
		Postcode: "SW1A 2AA",
	})
	fmt.Println(rsp, err)
}
