package example

import (
	"fmt"
	"os"

	"github.com/micro/services/clients/go/image"
)

// Delete an image previously uploaded.
func DeleteAnUploadedImage() {
	imageService := image.NewImageService(os.Getenv("MICRO_API_TOKEN"))
	rsp, err := imageService.Delete(&image.DeleteRequest{
		Url: "https://cdn.m3ocontent.com/micro/images/micro/41e23b39-48dd-42b6-9738-79a313414bb8/cat.png",
	})
	fmt.Println(rsp, err)
}
