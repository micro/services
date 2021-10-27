# Crypto

An [m3o.com](https://m3o.com) API. For example usage see [m3o.com/Crypto/api](https://m3o.com/Crypto/api).

Endpoints:

## History

Returns the history for the previous close


[https://m3o.com/crypto/api#History](https://m3o.com/crypto/api#History)

```go
package example

import(
	"fmt"
	"os"

	"github.com/micro/services/clients/go/crypto"
)

// Returns the history for the previous close
func GetPreviousClose() {
	cryptoService := crypto.NewCryptoService(os.Getenv("MICRO_API_TOKEN"))
	rsp, err := cryptoService.History(&crypto.HistoryRequest{
		Symbol: "BTCUSD",

	})
	fmt.Println(rsp, err)
}
```
## News

Get news related to a currency


[https://m3o.com/crypto/api#News](https://m3o.com/crypto/api#News)

```go
package example

import(
	"fmt"
	"os"

	"github.com/micro/services/clients/go/crypto"
)

// Get news related to a currency
func GetCryptocurrencyNews() {
	cryptoService := crypto.NewCryptoService(os.Getenv("MICRO_API_TOKEN"))
	rsp, err := cryptoService.News(&crypto.NewsRequest{
		Symbol: "BTCUSD",

	})
	fmt.Println(rsp, err)
}
```
## Price

Get the last price for a given crypto ticker


[https://m3o.com/crypto/api#Price](https://m3o.com/crypto/api#Price)

```go
package example

import(
	"fmt"
	"os"

	"github.com/micro/services/clients/go/crypto"
)

// Get the last price for a given crypto ticker
func GetCryptocurrencyPrice() {
	cryptoService := crypto.NewCryptoService(os.Getenv("MICRO_API_TOKEN"))
	rsp, err := cryptoService.Price(&crypto.PriceRequest{
		Symbol: "BTCUSD",

	})
	fmt.Println(rsp, err)
}
```
## Quote

Get the last quote for a given crypto ticker


[https://m3o.com/crypto/api#Quote](https://m3o.com/crypto/api#Quote)

```go
package example

import(
	"fmt"
	"os"

	"github.com/micro/services/clients/go/crypto"
)

// Get the last quote for a given crypto ticker
func GetAcryptocurrencyQuote() {
	cryptoService := crypto.NewCryptoService(os.Getenv("MICRO_API_TOKEN"))
	rsp, err := cryptoService.Quote(&crypto.QuoteRequest{
		Symbol: "BTCUSD",

	})
	fmt.Println(rsp, err)
}
```
