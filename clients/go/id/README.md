package id

import(
	"github.com/m3o/m3o-go/client"
)

func NewIdService(token string) *IdService {
	return &IdService{
		client: client.NewClient(&client.Options{
			Token: token,
		}),
	}
}

type IdService struct {
	client *client.Client
}


// Generate a unique ID. Defaults to uuid.
func (t *IdService) Generate(request *GenerateRequest) (*GenerateResponse, error) {
	rsp := &GenerateResponse{}
	return rsp, t.client.Call("id", "Generate", request, rsp)
}

// List the types of IDs available. No query params needed.
func (t *IdService) Types(request *TypesRequest) (*TypesResponse, error) {
	rsp := &TypesResponse{}
	return rsp, t.client.Call("id", "Types", request, rsp)
}




type GenerateRequest struct {
  // type of id e.g uuid, shortid, snowflake (64 bit), bigflake (128 bit)
  Type string `json:"type"`
}

type GenerateResponse struct {
  // the unique id generated
  Id string `json:"id"`
  // the type of id generated
  Type string `json:"type"`
}

type TypesRequest struct {
}

type TypesResponse struct {
  Types []string `json:"types"`
}

# { Id

An [m3o.com](https://m3o.com) API. For example usage see [m3o.com/Id/api](https://m3o.com/Id/api).

Endpoints:

#generate

// Generate a unique ID. Defaults to uuid.


[https://m3o.com/id/api#generate](https://m3o.com/id/api#generate)
#types

// List the types of IDs available. No query params needed.


[https://m3o.com/id/api#types](https://m3o.com/id/api#types)
