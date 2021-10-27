# Stock

An [m3o.com](https://m3o.com) API. For example usage see [m3o.com/Stock/api](https://m3o.com/Stock/api).

Endpoints:

## OrderBook

Get the historic order book and each trade by timestamp


[https://m3o.com/stock/api#OrderBook](https://m3o.com/stock/api#OrderBook)

```go
package example

import(
	"fmt"
	"os"

	"github.com/micro/services/clients/go/stock"
)

// Get the historic order book and each trade by timestamp
func OrderBookHistory() {
	stockService := stock.NewStockService(os.Getenv("MICRO_API_TOKEN"))
	rsp, err := stockService.OrderBook(&stock.OrderBookRequest{
		Date: "2020-10-01",
End: "2020-10-01T11:00:00Z",
Limit: 3,
Start: "2020-10-01T10:00:00Z",
Stock: "AAPL",

	})
	fmt.Println(rsp, err)
}
```
## Price

Get the last price for a given stock ticker


[https://m3o.com/stock/api#Price](https://m3o.com/stock/api#Price)

```go
package example

import(
	"fmt"
	"os"

	"github.com/micro/services/clients/go/stock"
)

// Get the last price for a given stock ticker
func GetAstockPrice() {
	stockService := stock.NewStockService(os.Getenv("MICRO_API_TOKEN"))
	rsp, err := stockService.Price(&stock.PriceRequest{
		Symbol: "AAPL",

	})
	fmt.Println(rsp, err)
}
```
## Quote

Get the last quote for the stock


[https://m3o.com/stock/api#Quote](https://m3o.com/stock/api#Quote)

```go
package example

import(
	"fmt"
	"os"

	"github.com/micro/services/clients/go/stock"
)

// Get the last quote for the stock
func GetAstockQuote() {
	stockService := stock.NewStockService(os.Getenv("MICRO_API_TOKEN"))
	rsp, err := stockService.Quote(&stock.QuoteRequest{
		Symbol: "AAPL",

	})
	fmt.Println(rsp, err)
}
```
## History

Get the historic open-close for a given day


[https://m3o.com/stock/api#History](https://m3o.com/stock/api#History)

```go
package example

import(
	"fmt"
	"os"

	"github.com/micro/services/clients/go/stock"
)

// Get the historic open-close for a given day
func GetHistoricData() {
	stockService := stock.NewStockService(os.Getenv("MICRO_API_TOKEN"))
	rsp, err := stockService.History(&stock.HistoryRequest{
		Date: "2020-10-01",
Stock: "AAPL",

	})
	fmt.Println(rsp, err)
}
```
