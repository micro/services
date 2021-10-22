package example

import (
	"fmt"
	"os"

	"github.com/micro/services/clients/go/user"
)

// Create a new user account. The email address and username for the account must be unique.
func CreateAnAccount() {
	userService := user.NewUserService(os.Getenv("MICRO_API_TOKEN"))
	rsp, err := userService.Create(&user.CreateRequest{
		Email:    "joe@example.com",
		Id:       "usrid-1",
		Password: "mySecretPass123",
		Username: "usrname-1",
	})
	fmt.Println(rsp, err)
}
