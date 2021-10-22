package example

import (
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
