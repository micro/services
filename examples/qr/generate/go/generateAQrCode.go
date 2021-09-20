package example

import (
	"fmt"
	"github.com/micro/services/clients/go/qr"
	"os"
)

//
func GenerateAqrCode() {
	qrService := qr.NewQrService(os.Getenv("MICRO_API_TOKEN"))
	rsp, err := qrService.Generate(&qr.GenerateRequest{
		Size: 300,
		Text: "https://m3o.com/qr",
	})
	fmt.Println(rsp, err)
}
