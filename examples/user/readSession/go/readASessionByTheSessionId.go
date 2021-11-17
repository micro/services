package example

import (
	"fmt"
	"os"

	"github.com/micro/services/clients/go/user"
)

// Read a session by the session id. In the event it has expired or is not found and error is returned.
func ReadAsessionByTheSessionId() {
	userService := user.NewUserService(os.Getenv("MICRO_API_TOKEN"))
	rsp, err := userService.ReadSession(&user.ReadSessionRequest{
		SessionId: "df91a612-5b24-4634-99ff-240220ab8f55",
	})
	fmt.Println(rsp, err)
}
