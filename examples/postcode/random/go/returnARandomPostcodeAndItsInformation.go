package example

import (
	"fmt"
	"github.com/micro/services/clients/go/postcode"
	"os"
)

// Return a random postcode and its related info
func ReturnArandomPostcodeAndItsInformation() {
	postcodeService := postcode.NewPostcodeService(os.Getenv("MICRO_API_TOKEN"))
	rsp, err := postcodeService.Random(&postcode.RandomRequest{})
	fmt.Println(rsp, err)
}
