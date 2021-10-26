package example

import (
	"fmt"
	"os"

	"github.com/micro/services/clients/go/stock"
)

// Get the historic order book and each trade by timestamp
func OrderBookHistory() {
	stockService := stock.NewStockService(os.Getenv("MICRO_API_TOKEN"))
	rsp, err := stockService.OrderBook(&stock.OrderBookRequest{
		Date:  "2020-10-01",
		End:   "2020-10-01T11:00:00Z",
		Limit: 3,
		Start: "2020-10-01T10:00:00Z",
		Stock: "AAPL",
	})
	fmt.Println(rsp, err)
}
