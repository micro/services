package example

import (
	"fmt"
	"os"

	"github.com/micro/services/clients/go/user"
)

// Update the account username or email
func UpdateAnAccount() {
	userService := user.NewUserService(os.Getenv("MICRO_API_TOKEN"))
	rsp, err := userService.Update(&user.UpdateRequest{
		Email: "joeotheremail@example.com",
		Id:    "usrid-1",
	})
	fmt.Println(rsp, err)
}
