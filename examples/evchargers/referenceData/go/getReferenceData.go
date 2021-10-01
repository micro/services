package example

import (
	"fmt"
	"github.com/micro/services/clients/go/evchargers"
	"os"
)

// Retrieve reference data as used by this API and in conjunction with the Search endpoint
func GetReferenceData() {
	evchargersService := evchargers.NewEvchargersService(os.Getenv("MICRO_API_TOKEN"))
	rsp, err := evchargersService.ReferenceData(&evchargers.ReferenceDataRequest{})
	fmt.Println(rsp, err)
}
