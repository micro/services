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
		Id: "8b98acbe-0b6a-4d66-a414-5ffbf666786f",
	})
	fmt.Println(rsp, err)
}
