package example

import (
	"fmt"
	"os"

	"github.com/micro/services/clients/go/image"
)

// Convert an image from one format (jpeg, png etc.) to an other either on the fly (from base64 to base64),
// or by uploading the conversion result.
func ConvertApngImageToAjpegTakenFromAurlAndSavedToAurlOnMicrosCdn() {
	imageService := image.NewImageService(os.Getenv("MICRO_API_TOKEN"))
	rsp, err := imageService.Convert(&image.ConvertRequest{
		Name: "cat.jpeg",
		Url:  "somewebsite.com/cat.png",
	})
	fmt.Println(rsp, err)
}
