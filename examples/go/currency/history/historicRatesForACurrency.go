package example

import (
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
