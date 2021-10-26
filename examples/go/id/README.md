# Id

An [m3o.com](https://m3o.com) API. For example usage see [m3o.com/Id/api](https://m3o.com/Id/api).

Endpoints:

## Generate

Generate a unique ID. Defaults to uuid.


[https://m3o.com/id/api#Generate](https://m3o.com/id/api#Generate)

```go
package example

import(
	"fmt"
	"os"

	"github.com/micro/services/clients/go/id"
)

// Generate a unique ID. Defaults to uuid.
func GenerateAuniqueId() {
	idService := id.NewIdService(os.Getenv("MICRO_API_TOKEN"))
	rsp, err := idService.Generate(&id.GenerateRequest{
		Type: "uuid",

	})
	fmt.Println(rsp, err)
}
```
## Generate

Generate a unique ID. Defaults to uuid.


[https://m3o.com/id/api#Generate](https://m3o.com/id/api#Generate)

```go
package example

import(
	"fmt"
	"os"

	"github.com/micro/services/clients/go/id"
)

// Generate a unique ID. Defaults to uuid.
func GenerateAshortId() {
	idService := id.NewIdService(os.Getenv("MICRO_API_TOKEN"))
	rsp, err := idService.Generate(&id.GenerateRequest{
		Type: "shortid",

	})
	fmt.Println(rsp, err)
}
```
## Generate

Generate a unique ID. Defaults to uuid.


[https://m3o.com/id/api#Generate](https://m3o.com/id/api#Generate)

```go
package example

import(
	"fmt"
	"os"

	"github.com/micro/services/clients/go/id"
)

// Generate a unique ID. Defaults to uuid.
func GenerateAsnowflakeId() {
	idService := id.NewIdService(os.Getenv("MICRO_API_TOKEN"))
	rsp, err := idService.Generate(&id.GenerateRequest{
		Type: "snowflake",

	})
	fmt.Println(rsp, err)
}
```
## Generate

Generate a unique ID. Defaults to uuid.


[https://m3o.com/id/api#Generate](https://m3o.com/id/api#Generate)

```go
package example

import(
	"fmt"
	"os"

	"github.com/micro/services/clients/go/id"
)

// Generate a unique ID. Defaults to uuid.
func GenerateAbigflakeId() {
	idService := id.NewIdService(os.Getenv("MICRO_API_TOKEN"))
	rsp, err := idService.Generate(&id.GenerateRequest{
		Type: "bigflake",

	})
	fmt.Println(rsp, err)
}
```
## Types

List the types of IDs available. No query params needed.


[https://m3o.com/id/api#Types](https://m3o.com/id/api#Types)

```go
package example

import(
	"fmt"
	"os"

	"github.com/micro/services/clients/go/id"
)

// List the types of IDs available. No query params needed.
func ListTheTypesOfIdsAvailable() {
	idService := id.NewIdService(os.Getenv("MICRO_API_TOKEN"))
	rsp, err := idService.Types(&id.TypesRequest{
		
	})
	fmt.Println(rsp, err)
}
```
