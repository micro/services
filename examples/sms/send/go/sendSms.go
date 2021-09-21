package example

import (
	"fmt"
	"github.com/micro/services/clients/go/sms"
	"os"
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
