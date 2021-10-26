# Currency

An [m3o.com](https://m3o.com) API. For example usage see [m3o.com/Currency/api](https://m3o.com/Currency/api).

Endpoints:

## Codes

Codes returns the supported currency codes for the API


[https://m3o.com/currency/api#Codes](https://m3o.com/currency/api#Codes)

```go
package example

import(
	"fmt"
	"os"

	"github.com/micro/services/clients/go/currency"
)

// Codes returns the supported currency codes for the API
func GetSupportedCodes() {
	currencyService := currency.NewCurrencyService(os.Getenv("MICRO_API_TOKEN"))
	rsp, err := currencyService.Codes(&currency.CodesRequest{
		
	})
	fmt.Println(rsp, err)
}
```
## Rates

Rates returns the currency rates for a given code e.g USD


[https://m3o.com/currency/api#Rates](https://m3o.com/currency/api#Rates)

```go
package example

import(
	"fmt"
	"os"

	"github.com/micro/services/clients/go/currency"
)

// Rates returns the currency rates for a given code e.g USD
func GetRatesForUsd() {
	currencyService := currency.NewCurrencyService(os.Getenv("MICRO_API_TOKEN"))
	rsp, err := currencyService.Rates(&currency.RatesRequest{
		Code: "USD",

	})
	fmt.Println(rsp, err)
}
```
## Convert

Convert returns the currency conversion rate between two pairs e.g USD/GBP


[https://m3o.com/currency/api#Convert](https://m3o.com/currency/api#Convert)

```go
package example

import(
	"fmt"
	"os"

	"github.com/micro/services/clients/go/currency"
)

// Convert returns the currency conversion rate between two pairs e.g USD/GBP
func ConvertUsdToGbp() {
	currencyService := currency.NewCurrencyService(os.Getenv("MICRO_API_TOKEN"))
	rsp, err := currencyService.Convert(&currency.ConvertRequest{
		From: "USD",
To: "GBP",

	})
	fmt.Println(rsp, err)
}
```
## Convert

Convert returns the currency conversion rate between two pairs e.g USD/GBP


[https://m3o.com/currency/api#Convert](https://m3o.com/currency/api#Convert)

```go
package example

import(
	"fmt"
	"os"

	"github.com/micro/services/clients/go/currency"
)

// Convert returns the currency conversion rate between two pairs e.g USD/GBP
func Convert10usdToGbp() {
	currencyService := currency.NewCurrencyService(os.Getenv("MICRO_API_TOKEN"))
	rsp, err := currencyService.Convert(&currency.ConvertRequest{
		Amount: 10,
From: "USD",
To: "GBP",

	})
	fmt.Println(rsp, err)
}
```
## History

Returns the historic rates for a currency on a given date


[https://m3o.com/currency/api#History](https://m3o.com/currency/api#History)

```go
package example

import(
	"fmt"
	"os"

	"github.com/micro/services/clients/go/currency"
)

// Returns the historic rates for a currency on a given date
func HistoricRatesForAcurrency() {
	currencyService := currency.NewCurrencyService(os.Getenv("MICRO_API_TOKEN"))
	rsp, err := currencyService.History(&currency.HistoryRequest{
		Code: "USD",
Date: "2021-05-30",

	})
	fmt.Println(rsp, err)
}
```
