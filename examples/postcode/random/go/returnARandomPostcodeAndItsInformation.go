package example

import (
	"fmt"
	"os"

	"github.com/micro/services/clients/go/postcode"
)

// Return a random postcode and its related info
func ReturnArandomPostcodeAndItsInformation() {
	postcodeService := postcode.NewPostcodeService(os.Getenv("MICRO_API_TOKEN"))
	rsp, err := postcodeService.Random(&postcode.RandomRequest{})
	fmt.Println(rsp, err)
}
