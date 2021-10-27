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
		ConfirmPassword: "myEvenMoreSecretPass123",
		NewPassword:     "myEvenMoreSecretPass123",
		OldPassword:     "mySecretPass123",
	})
	fmt.Println(rsp, err)
}
