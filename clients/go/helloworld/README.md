# Helloworld

An [m3o.com](https://m3o.com) API. For example usage see [m3o.com/Helloworld/api](https://m3o.com/Helloworld/api).

Endpoints:

## Call

Call returns a personalised "Hello $name" response


[https://m3o.com/helloworld/api#Call](https://m3o.com/helloworld/api#Call)

```go
package example

import(
	"fmt"
	"os"

	"github.com/micro/services/clients/go/helloworld"
)

// Call returns a personalised "Hello $name" response
func CallTheHelloworldService() {
	helloworldService := helloworld.NewHelloworldService(os.Getenv("MICRO_API_TOKEN"))
	rsp, err := helloworldService.Call(&helloworld.CallRequest{
		Name: "John",

	})
	fmt.Println(rsp, err)
}
```
## Stream

Stream returns a stream of "Hello $name" responses


[https://m3o.com/helloworld/api#Stream](https://m3o.com/helloworld/api#Stream)

```go
package example

import(
	"fmt"
	"os"

	"github.com/micro/services/clients/go/helloworld"
)

// Stream returns a stream of "Hello $name" responses
func StreamsAreCurrentlyTemporarilyNotSupportedInClients() {
	helloworldService := helloworld.NewHelloworldService(os.Getenv("MICRO_API_TOKEN"))
	rsp, err := helloworldService.Stream(&helloworld.StreamRequest{
		Name: "not supported",

	})
	fmt.Println(rsp, err)
}
```
