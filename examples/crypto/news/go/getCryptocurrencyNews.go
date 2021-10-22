package example

import (
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
