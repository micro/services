package example

import (
	"fmt"
	"os"

	"github.com/micro/services/clients/go/holidays"
)

// Get the list of countries that are supported by this API
func ListCountries() {
	holidaysService := holidays.NewHolidaysService(os.Getenv("MICRO_API_TOKEN"))
	rsp, err := holidaysService.Countries(&holidays.CountriesRequest{})
	fmt.Println(rsp, err)
}
