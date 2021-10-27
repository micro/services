package example

import (
	"fmt"
	"os"

	"github.com/micro/services/clients/go/image"
)

// Upload an image by either sending a base64 encoded image to this endpoint or a URL.
// To resize an image before uploading, see the Resize endpoint.
func UploadAnImageFromAurlToMicrosCdn() {
	imageService := image.NewImageService(os.Getenv("MICRO_API_TOKEN"))
	rsp, err := imageService.Upload(&image.UploadRequest{
		Name: "cat.jpeg",
		Url:  "somewebsite.com/cat.png",
	})
	fmt.Println(rsp, err)
}
