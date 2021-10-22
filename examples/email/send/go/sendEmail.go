package example

import (
	"fmt"
	"os"

	"github.com/micro/services/clients/go/email"
)

// Send an email by passing in from, to, subject, and a text or html body
func SendEmail() {
	emailService := email.NewEmailService(os.Getenv("MICRO_API_TOKEN"))
	rsp, err := emailService.Send(&email.SendRequest{
		From:    "Awesome Dot Com",
		Subject: "Email verification",
		TextBody: `Hi there,

Please verify your email by clicking this link: $micro_verification_link`,
	})
	fmt.Println(rsp, err)
}
