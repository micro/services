package example

import (
	"fmt"
	"os"

	"github.com/micro/services/clients/go/geocoding"
)

// Lookup returns a geocoded address including normalized address and gps coordinates. All fields are optional, provide more to get more accurate results
func GeocodeAnAddress() {
	geocodingService := geocoding.NewGeocodingService(os.Getenv("MICRO_API_TOKEN"))
	rsp, err := geocodingService.Lookup(&geocoding.LookupRequest{
		Address:  "10 russell st",
		City:     "london",
		Country:  "uk",
		Postcode: "wc2b",
	})
	fmt.Println(rsp, err)
}
