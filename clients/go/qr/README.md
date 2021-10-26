# Qr

An [m3o.com](https://m3o.com) API. For example usage see [m3o.com/Qr/api](https://m3o.com/Qr/api).

Endpoints:

## Generate

Generate a QR code with a specific text and size


[https://m3o.com/qr/api#Generate](https://m3o.com/qr/api#Generate)

```go
package example

import(
	"fmt"
	"os"

	"github.com/micro/services/clients/go/qr"
)

// Generate a QR code with a specific text and size
func GenerateAqrCode() {
	qrService := qr.NewQrService(os.Getenv("MICRO_API_TOKEN"))
	rsp, err := qrService.Generate(&qr.GenerateRequest{
		Size: 300,
Text: "https://m3o.com/qr",

	})
	fmt.Println(rsp, err)
}
```
