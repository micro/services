package example

import (
	"fmt"
	"github.com/micro/services/clients/go/prayer"
	"os"
)

// Get the prayer (salah) times for a location on a given date
func PrayerTimes() {
	prayerService := prayer.NewPrayerService(os.Getenv("MICRO_API_TOKEN"))
	rsp, err := prayerService.Times(&prayer.TimesRequest{
		Location: "london",
	})
	fmt.Println(rsp, err)
}
