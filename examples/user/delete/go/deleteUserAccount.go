package example

import (
	"fmt"
	"os"

	"github.com/micro/services/clients/go/user"
)

// Delete an account by id
func DeleteUserAccount() {
	userService := user.NewUserService(os.Getenv("MICRO_API_TOKEN"))
	rsp, err := userService.Delete(&user.DeleteRequest{
		Id: "fdf34f34f34-f34f34-f43f43f34-f4f34f",
	})
	fmt.Println(rsp, err)
}
