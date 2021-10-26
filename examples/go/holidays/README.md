# Holidays

An [m3o.com](https://m3o.com) API. For example usage see [m3o.com/Holidays/api](https://m3o.com/Holidays/api).

Endpoints:

## Countries

Get the list of countries that are supported by this API


[https://m3o.com/holidays/api#Countries](https://m3o.com/holidays/api#Countries)

```go
package example

import(
	"fmt"
	"os"

	"github.com/micro/services/clients/go/holidays"
)

// Get the list of countries that are supported by this API
func ListCountries() {
	holidaysService := holidays.NewHolidaysService(os.Getenv("MICRO_API_TOKEN"))
	rsp, err := holidaysService.Countries(&holidays.CountriesRequest{
		
	})
	fmt.Println(rsp, err)
}
```
## List

List the holiday dates for a given country and year


[https://m3o.com/holidays/api#List](https://m3o.com/holidays/api#List)

```go
package example

import(
	"fmt"
	"os"

	"github.com/micro/services/clients/go/holidays"
)

// List the holiday dates for a given country and year
func GetHolidays() {
	holidaysService := holidays.NewHolidaysService(os.Getenv("MICRO_API_TOKEN"))
	rsp, err := holidaysService.List(&holidays.ListRequest{
		Year: 2022,

	})
	fmt.Println(rsp, err)
}
```
