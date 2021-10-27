package example

import (
	"fmt"
	"os"

	"github.com/micro/services/clients/go/user"
)

// Verify the email address of an account from a token sent in an email to the user.
func VerifyEmail() {
	userService := user.NewUserService(os.Getenv("MICRO_API_TOKEN"))
	rsp, err := userService.VerifyEmail(&user.VerifyEmailRequest{
		Token: "t2323t232t",
	})
	fmt.Println(rsp, err)
}
