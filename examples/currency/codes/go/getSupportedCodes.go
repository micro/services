package example

import (
	"fmt"
	"os"

	"github.com/micro/services/clients/go/currency"
)

// Codes returns the supported currency codes for the API
func GetSupportedCodes() {
	currencyService := currency.NewCurrencyService(os.Getenv("MICRO_API_TOKEN"))
	rsp, err := currencyService.Codes(&currency.CodesRequest{})
	fmt.Println(rsp, err)
}
