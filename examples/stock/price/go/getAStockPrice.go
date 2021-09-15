package example

import (
	"fmt"
	"github.com/micro/services/clients/go/stock"
	"os"
)

// Get the last price for a given stock ticker
func GetAstockPrice() {
	stockService := stock.NewStockService(os.Getenv("MICRO_API_TOKEN"))
	rsp, err := stockService.Price(&stock.PriceRequest{
		Symbol: "AAPL",
	})
	fmt.Println(rsp, err)
}
