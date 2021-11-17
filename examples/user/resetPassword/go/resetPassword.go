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
		Code:            "012345",
		ConfirmPassword: "NewPassword1",
		Email:           "joe@example.com",
		NewPassword:     "NewPassword1",
	})
	fmt.Println(rsp, err)
}
