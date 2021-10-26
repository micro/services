# Evchargers

An [m3o.com](https://m3o.com) API. For example usage see [m3o.com/Evchargers/api](https://m3o.com/Evchargers/api).

Endpoints:

## Search

Search by giving a coordinate and a max distance, or bounding box and optional filters


[https://m3o.com/evchargers/api#Search](https://m3o.com/evchargers/api#Search)

```go
package example

import(
	"fmt"
	"os"

	"github.com/micro/services/clients/go/evchargers"
)

// Search by giving a coordinate and a max distance, or bounding box and optional filters
func SearchByLocation() {
	evchargersService := evchargers.NewEvchargersService(os.Getenv("MICRO_API_TOKEN"))
	rsp, err := evchargersService.Search(&evchargers.SearchRequest{
		Distance: 2000,
Location: &evchargers.Coordinates{
	Latitude: 51.53336351319885,
	Longitude: -0.0252,
},

	})
	fmt.Println(rsp, err)
}
```
## Search

Search by giving a coordinate and a max distance, or bounding box and optional filters


[https://m3o.com/evchargers/api#Search](https://m3o.com/evchargers/api#Search)

```go
package example

import(
	"fmt"
	"os"

	"github.com/micro/services/clients/go/evchargers"
)

// Search by giving a coordinate and a max distance, or bounding box and optional filters
func SearchByBoundingBox() {
	evchargersService := evchargers.NewEvchargersService(os.Getenv("MICRO_API_TOKEN"))
	rsp, err := evchargersService.Search(&evchargers.SearchRequest{
		Box: &evchargers.BoundingBox{
		},

	})
	fmt.Println(rsp, err)
}
```
## Search

Search by giving a coordinate and a max distance, or bounding box and optional filters


[https://m3o.com/evchargers/api#Search](https://m3o.com/evchargers/api#Search)

```go
package example

import(
	"fmt"
	"os"

	"github.com/micro/services/clients/go/evchargers"
)

// Search by giving a coordinate and a max distance, or bounding box and optional filters
func SearchWithFiltersFastChargersOnly() {
	evchargersService := evchargers.NewEvchargersService(os.Getenv("MICRO_API_TOKEN"))
	rsp, err := evchargersService.Search(&evchargers.SearchRequest{
		Distance: 2000,
Levels: []string{"3"},
Location: &evchargers.Coordinates{
	Latitude: 51.53336351319885,
	Longitude: -0.0252,
},

	})
	fmt.Println(rsp, err)
}
```
## ReferenceData

Retrieve reference data as used by this API and in conjunction with the Search endpoint


[https://m3o.com/evchargers/api#ReferenceData](https://m3o.com/evchargers/api#ReferenceData)

```go
package example

import(
	"fmt"
	"os"

	"github.com/micro/services/clients/go/evchargers"
)

// Retrieve reference data as used by this API and in conjunction with the Search endpoint
func GetReferenceData() {
	evchargersService := evchargers.NewEvchargersService(os.Getenv("MICRO_API_TOKEN"))
	rsp, err := evchargersService.ReferenceData(&evchargers.ReferenceDataRequest{
		
	})
	fmt.Println(rsp, err)
}
```
