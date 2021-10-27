# Image

An [m3o.com](https://m3o.com) API. For example usage see [m3o.com/Image/api](https://m3o.com/Image/api).

Endpoints:

## Upload

Upload an image by either sending a base64 encoded image to this endpoint or a URL.
To resize an image before uploading, see the Resize endpoint.


[https://m3o.com/image/api#Upload](https://m3o.com/image/api#Upload)

```go
package example

import(
	"fmt"
	"os"

	"github.com/micro/services/clients/go/image"
)

// Upload an image by either sending a base64 encoded image to this endpoint or a URL.
// To resize an image before uploading, see the Resize endpoint.
func UploadAbase64imageToMicrosCdn() {
	imageService := image.NewImageService(os.Getenv("MICRO_API_TOKEN"))
	rsp, err := imageService.Upload(&image.UploadRequest{
		Base64: "data:image/png;base64, iVBORw0KGgoAAAANSUhEUgAAADIAAAAyCAYAAAAeP4ixAAAAx0lEQVR4nOzaMaoDMQyE4ZHj+x82vVdhwQoTkzKQEcwP5r0ihT7sbjUTeAJ4HCegXQJYfOYefOyjDuBiz3yjwJBoCIl6QZOeUjTC1Ix1IxEJXF9+0KWsf2bD4bn37OO/c/wuQ9QyRC1D1DJELUPUMkQtQ9QyRC1D1DJELUPUMkQtQ9QyRC1D1DJELUPUMkQtQ9Sa/NG94Tf3j4WBdaxudMEkn4IM2rZBA0wBrvo7aOcpj2emXvLeVt0IGm0GVXUj91mvAAAA//+V2CZl+4AKXwAAAABJRU5ErkJggg==",
Name: "cat.jpeg",

	})
	fmt.Println(rsp, err)
}
```
## Upload

Upload an image by either sending a base64 encoded image to this endpoint or a URL.
To resize an image before uploading, see the Resize endpoint.


[https://m3o.com/image/api#Upload](https://m3o.com/image/api#Upload)

```go
package example

import(
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
Url: "somewebsite.com/cat.png",

	})
	fmt.Println(rsp, err)
}
```
## Resize

Resize an image on the fly without storing it (by sending and receiving a base64 encoded image), or resize and upload depending on parameters.
If one of width or height is 0, the image aspect ratio is preserved.
Optional cropping.


[https://m3o.com/image/api#Resize](https://m3o.com/image/api#Resize)

```go
package example

import(
	"fmt"
	"os"

	"github.com/micro/services/clients/go/image"
)

// Resize an image on the fly without storing it (by sending and receiving a base64 encoded image), or resize and upload depending on parameters.
// If one of width or height is 0, the image aspect ratio is preserved.
// Optional cropping.
func Base64toHostedImage() {
	imageService := image.NewImageService(os.Getenv("MICRO_API_TOKEN"))
	rsp, err := imageService.Resize(&image.ResizeRequest{
		Base64: "data:image/png;base64, iVBORw0KGgoAAAANSUhEUgAAAAUAAAAFCAYAAACNbyblAAAAHElEQVQI12P4//8/w38GIAXDIBKE0DHxgljNBAAO9TXL0Y4OHwAAAABJRU5ErkJggg==",
Height: 100,
Name: "cat.png",
Width: 100,

	})
	fmt.Println(rsp, err)
}
```
## Resize

Resize an image on the fly without storing it (by sending and receiving a base64 encoded image), or resize and upload depending on parameters.
If one of width or height is 0, the image aspect ratio is preserved.
Optional cropping.


[https://m3o.com/image/api#Resize](https://m3o.com/image/api#Resize)

```go
package example

import(
	"fmt"
	"os"

	"github.com/micro/services/clients/go/image"
)

// Resize an image on the fly without storing it (by sending and receiving a base64 encoded image), or resize and upload depending on parameters.
// If one of width or height is 0, the image aspect ratio is preserved.
// Optional cropping.
func Base64toBase64image() {
	imageService := image.NewImageService(os.Getenv("MICRO_API_TOKEN"))
	rsp, err := imageService.Resize(&image.ResizeRequest{
		Base64: "data:image/png;base64, iVBORw0KGgoAAAANSUhEUgAAAAUAAAAFCAYAAACNbyblAAAAHElEQVQI12P4//8/w38GIAXDIBKE0DHxgljNBAAO9TXL0Y4OHwAAAABJRU5ErkJggg==",
Height: 100,
Width: 100,

	})
	fmt.Println(rsp, err)
}
```
## Resize

Resize an image on the fly without storing it (by sending and receiving a base64 encoded image), or resize and upload depending on parameters.
If one of width or height is 0, the image aspect ratio is preserved.
Optional cropping.


[https://m3o.com/image/api#Resize](https://m3o.com/image/api#Resize)

```go
package example

import(
	"fmt"
	"os"

	"github.com/micro/services/clients/go/image"
)

// Resize an image on the fly without storing it (by sending and receiving a base64 encoded image), or resize and upload depending on parameters.
// If one of width or height is 0, the image aspect ratio is preserved.
// Optional cropping.
func Base64toBase64imageWithCropping() {
	imageService := image.NewImageService(os.Getenv("MICRO_API_TOKEN"))
	rsp, err := imageService.Resize(&image.ResizeRequest{
		Base64: "data:image/png;base64, iVBORw0KGgoAAAANSUhEUgAAAAUAAAAFCAYAAACNbyblAAAAHElEQVQI12P4//8/w38GIAXDIBKE0DHxgljNBAAO9TXL0Y4OHwAAAABJRU5ErkJggg==",
CropOptions: &image.CropOptions{
		Height: 50,
	Width: 50,
},
Height: 100,
Width: 100,

	})
	fmt.Println(rsp, err)
}
```
## Convert

Convert an image from one format (jpeg, png etc.) to an other either on the fly (from base64 to base64),
or by uploading the conversion result.


[https://m3o.com/image/api#Convert](https://m3o.com/image/api#Convert)

```go
package example

import(
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
Url: "somewebsite.com/cat.png",

	})
	fmt.Println(rsp, err)
}
```
