# Vehicle

An [m3o.com](https://m3o.com) API. For example usage see [m3o.com/Vehicle/api](https://m3o.com/Vehicle/api).

Endpoints:

## Lookup

Lookup a UK vehicle by it's registration number


[https://m3o.com/vehicle/api#Lookup](https://m3o.com/vehicle/api#Lookup)

```go
package example

import(
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
```
