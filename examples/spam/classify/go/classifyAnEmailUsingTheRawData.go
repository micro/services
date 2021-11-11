package example

import (
	"fmt"
	"os"

	"github.com/micro/services/clients/go/spam"
)

// Check whether an email is likely to be spam based on its attributes
func ClassifyAnEmailUsingTheRawData() {
	spamService := spam.NewSpamService(os.Getenv("MICRO_API_TOKEN"))
	rsp, err := spamService.Classify(&spam.ClassifyRequest{})
	fmt.Println(rsp, err)
}
