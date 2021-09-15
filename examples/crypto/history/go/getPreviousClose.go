package example

import (
	"fmt"
	"github.com/micro/services/clients/go/crypto"
	"os"
)

// Returns the history for the previous close
func GetPreviousClose() {
	cryptoService := crypto.NewCryptoService(os.Getenv("MICRO_API_TOKEN"))
	rsp, err := cryptoService.History(&crypto.HistoryRequest{
		Symbol: "BTCUSD",
	})
	fmt.Println(rsp, err)
}
