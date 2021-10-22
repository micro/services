package example

import (
	"fmt"
	"os"

	"github.com/micro/services/clients/go/user"
)

// Read an account by id, username or email. Only one need to be specified.
func ReadAccountByUsernameOrEmail() {
	userService := user.NewUserService(os.Getenv("MICRO_API_TOKEN"))
	rsp, err := userService.Read(&user.ReadRequest{
		Username: "usrname-1",
	})
	fmt.Println(rsp, err)
}
