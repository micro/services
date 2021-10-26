package example

import (
	"fmt"
	"os"

	"github.com/micro/services/clients/go/forex"
)

// Returns the data for the previous close
func GetPreviousClose() {
	forexService := forex.NewForexService(os.Getenv("MICRO_API_TOKEN"))
	rsp, err := forexService.History(&forex.HistoryRequest{
		Symbol: "GBPUSD",
	})
	fmt.Println(rsp, err)
}
