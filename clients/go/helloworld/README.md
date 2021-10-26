package helloworld

import(
	"github.com/m3o/m3o-go/client"
)

func NewHelloworldService(token string) *HelloworldService {
	return &HelloworldService{
		client: client.NewClient(&client.Options{
			Token: token,
		}),
	}
}

type HelloworldService struct {
	client *client.Client
}


// Call returns a personalised "Hello $name" response
func (t *HelloworldService) Call(request *CallRequest) (*CallResponse, error) {
	rsp := &CallResponse{}
	return rsp, t.client.Call("helloworld", "Call", request, rsp)
}

// Stream returns a stream of "Hello $name" responses
func (t *HelloworldService) Stream(request *StreamRequest) (*StreamResponse, error) {
	rsp := &StreamResponse{}
	return rsp, t.client.Call("helloworld", "Stream", request, rsp)
}




type CallRequest struct {
  Name string `json:"name"`
}

type CallResponse struct {
  Message string `json:"message"`
}

type StreamRequest struct {
  // the number of messages to send back
  Messages int64 `json:"messages,string"`
  Name string `json:"name"`
}

type StreamResponse struct {
  Message string `json:"message"`
}

# { Helloworld

An [m3o.com](https://m3o.com) API. For example usage see [m3o.com/Helloworld/api](https://m3o.com/Helloworld/api).

Endpoints:

#call

// Call returns a personalised "Hello $name" response


[https://m3o.com/helloworld/api#call](https://m3o.com/helloworld/api#call)
#stream

// Stream returns a stream of "Hello $name" responses


[https://m3o.com/helloworld/api#stream](https://m3o.com/helloworld/api#stream)
