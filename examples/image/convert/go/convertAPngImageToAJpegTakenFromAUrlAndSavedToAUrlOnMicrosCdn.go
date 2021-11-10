package example

import (
	"fmt"
	"os"

	"github.com/micro/services/clients/go/image"
)

// Convert an image from one format (jpeg, png etc.) to an other either on the fly (from base64 to base64),
// or by uploading the conversion result.
// To use the file parameter you need to send the request as a multipart/form-data rather than the usual application/json
// with each parameter as a form field.
func ConvertApngImageToAjpegTakenFromAurlAndSavedToAurlOnMicrosCdn() {
	imageService := image.NewImageService(os.Getenv("MICRO_API_TOKEN"))
	rsp, err := imageService.Convert(&image.ConvertRequest{
		Name: "cat.jpeg",
		Url:  "somewebsite.com/cat.png",
	})
	fmt.Println(rsp, err)
}
