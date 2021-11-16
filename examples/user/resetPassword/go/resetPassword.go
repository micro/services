package example

import (
	"fmt"
	"os"

	"github.com/micro/services/clients/go/user"
)

// Reset password with the code sent by the "SendPasswordResetEmail" endoint.
func ResetPassword() {
	userService := user.NewUserService(os.Getenv("MICRO_API_TOKEN"))
	rsp, err := userService.ResetPassword(&user.ResetPasswordRequest{
		Code:            "some-code-from-email",
		ConfirmPassword: "newpass123",
		NewPassword:     "newpass123",
	})
	fmt.Println(rsp, err)
}
