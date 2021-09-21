package example

import (
	"fmt"
	"github.com/micro/services/clients/go/forex"
	"os"
)

// Get the latest price for a given forex ticker
func GetAnFxPrice() {
	forexService := forex.NewForexService(os.Getenv("MICRO_API_TOKEN"))
	rsp, err := forexService.Price(&forex.PriceRequest{
		Symbol: "GBPUSD",
	})
	fmt.Println(rsp, err)
}
