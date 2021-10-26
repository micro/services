# Stream

An [m3o.com](https://m3o.com) API. For example usage see [m3o.com/Stream/api](https://m3o.com/Stream/api).

Endpoints:

## Publish

Publish a message to the stream. Specify a topic to group messages for a specific topic.


[https://m3o.com/stream/api#Publish](https://m3o.com/stream/api#Publish)

```go
package example

import(
	"fmt"
	"os"

	"github.com/micro/services/clients/go/stream"
)

// Publish a message to the stream. Specify a topic to group messages for a specific topic.
func PublishAmessage() {
	streamService := stream.NewStreamService(os.Getenv("MICRO_API_TOKEN"))
	rsp, err := streamService.Publish(&stream.PublishRequest{
		Message: map[string]interface{}{
	"id": "1",
	"type": "signup",
	"user": "john",
},
Topic: "events",

	})
	fmt.Println(rsp, err)
}
```
## Subscribe

Subscribe to messages for a given topic.


[https://m3o.com/stream/api#Subscribe](https://m3o.com/stream/api#Subscribe)

```go
package example

import(
	"fmt"
	"os"

	"github.com/micro/services/clients/go/stream"
)

// Subscribe to messages for a given topic.
func SubscribeToAtopic() {
	streamService := stream.NewStreamService(os.Getenv("MICRO_API_TOKEN"))
	rsp, err := streamService.Subscribe(&stream.SubscribeRequest{
		Topic: "events",

	})
	fmt.Println(rsp, err)
}
```
