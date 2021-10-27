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
		SessionId: "sds34s34s34-s34s34-s43s43s34-s4s34s",
	})
	fmt.Println(rsp, err)
}
