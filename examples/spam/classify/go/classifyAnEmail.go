package example

import (
	"fmt"
	"os"

	"github.com/micro/services/clients/go/spam"
)

// Check whether an email is likely to be spam based on its attributes
func ClassifyAnEmail() {
	spamService := spam.NewSpamService(os.Getenv("MICRO_API_TOKEN"))
	rsp, err := spamService.Classify(&spam.ClassifyRequest{
		From:    "noreply@m3o.com",
		Subject: "Welcome",
		To:      "hello@example.com",
	})
	fmt.Println(rsp, err)
}
