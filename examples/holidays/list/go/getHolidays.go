package example

import (
	"fmt"
	"github.com/micro/services/clients/go/holidays"
	"os"
)

// List the holiday dates for a given country and year
func GetHolidays() {
	holidaysService := holidays.NewHolidaysService(os.Getenv("MICRO_API_TOKEN"))
	rsp, err := holidaysService.List(&holidays.ListRequest{
		Year: 2022,
	})
	fmt.Println(rsp, err)
}
