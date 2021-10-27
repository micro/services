package example

import (
	"fmt"
	"os"

	"github.com/micro/services/clients/go/sms"
)

// Send an SMS.
func SendSms() {
	smsService := sms.NewSmsService(os.Getenv("MICRO_API_TOKEN"))
	rsp, err := smsService.Send(&sms.SendRequest{
		From:    "Alice",
		Message: "Hi there!",
		To:      "+447681129",
	})
	fmt.Println(rsp, err)
}
