package example

import (
	"fmt"
	"github.com/micro/services/clients/go/stock"
	"os"
)

// Get the last quote for the stock
func GetAstockQuote() {
	stockService := stock.NewStockService(os.Getenv("MICRO_API_TOKEN"))
	rsp, err := stockService.Quote(&stock.QuoteRequest{
		Symbol: "AAPL",
	})
	fmt.Println(rsp, err)
}
