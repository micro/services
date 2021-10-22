package example

import (
	"fmt"
	"os"

	"github.com/micro/services/clients/go/vehicle"
)

// Lookup a UK vehicle by it's registration number
func LookupVehicle() {
	vehicleService := vehicle.NewVehicleService(os.Getenv("MICRO_API_TOKEN"))
	rsp, err := vehicleService.Lookup(&vehicle.LookupRequest{
		Registration: "LC60OTA",
	})
	fmt.Println(rsp, err)
}
