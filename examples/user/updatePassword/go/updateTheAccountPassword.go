package example

import (
	"fmt"
	"os"

	"github.com/micro/services/clients/go/user"
)

// Update the account password
func UpdateTheAccountPassword() {
	userService := user.NewUserService(os.Getenv("MICRO_API_TOKEN"))
	rsp, err := userService.UpdatePassword(&user.UpdatePasswordRequest{
		ConfirmPassword: "Password2",
		NewPassword:     "Password2",
		OldPassword:     "Password1",
	})
	fmt.Println(rsp, err)
}
