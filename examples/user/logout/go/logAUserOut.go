package example

import (
	"fmt"
	"os"

	"github.com/micro/services/clients/go/user"
)

// Logout a user account
func LogAuserOut() {
	userService := user.NewUserService(os.Getenv("MICRO_API_TOKEN"))
	rsp, err := userService.Logout(&user.LogoutRequest{
		SessionId: "df91a612-5b24-4634-99ff-240220ab8f55",
	})
	fmt.Println(rsp, err)
}
