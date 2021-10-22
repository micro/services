package example

import (
	"fmt"
	"os"

	"github.com/micro/services/clients/go/user"
)

// Login using username or email. The response will return a new session for successful login,
// 401 in the case of login failure and 500 for any other error
func LogAuserIn() {
	userService := user.NewUserService(os.Getenv("MICRO_API_TOKEN"))
	rsp, err := userService.Login(&user.LoginRequest{
		Email:    "joe@example.com",
		Password: "mySecretPass123",
	})
	fmt.Println(rsp, err)
}
