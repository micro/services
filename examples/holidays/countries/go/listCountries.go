package example

import (
	"fmt"
	"github.com/micro/services/clients/go/holidays"
	"os"
)

//
func ListCountries() {
	holidaysService := holidays.NewHolidaysService(os.Getenv("MICRO_API_TOKEN"))
	rsp, err := holidaysService.Countries(&holidays.CountriesRequest{})
	fmt.Println(rsp, err)
}
