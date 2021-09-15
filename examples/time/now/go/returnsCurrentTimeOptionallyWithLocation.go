package example

import (
	"fmt"
	"github.com/micro/services/clients/go/time"
	"os"
)

// Get the current time
func ReturnsCurrentTimeOptionallyWithLocation() {
	timeService := time.NewTimeService(os.Getenv("MICRO_API_TOKEN"))
	rsp, err := timeService.Now(&time.NowRequest{})
	fmt.Println(rsp, err)
}
