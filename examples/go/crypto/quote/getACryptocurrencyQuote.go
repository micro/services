package example

import (
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
