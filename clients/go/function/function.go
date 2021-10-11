package function

import (
	"github.com/m3o/m3o-go/client"
)

func NewFunctionService(token string) *FunctionService {
	return &FunctionService{
		client: client.NewClient(&client.Options{
			Token: token,
		}),
	}
}

type FunctionService struct {
	client *client.Client
}

// Call a function
func (t *FunctionService) Call(request *CallRequest) (*CallResponse, error) {
	rsp := &CallResponse{}
	return rsp, t.client.Call("function", "Call", request, rsp)
}

//
func (t *FunctionService) Delete(request *DeleteRequest) (*DeleteResponse, error) {
	rsp := &DeleteResponse{}
	return rsp, t.client.Call("function", "Delete", request, rsp)
}

// Deploy a group of functions
func (t *FunctionService) Deploy(request *DeployRequest) (*DeployResponse, error) {
	rsp := &DeployResponse{}
	return rsp, t.client.Call("function", "Deploy", request, rsp)
}

//
func (t *FunctionService) List(request *ListRequest) (*ListResponse, error) {
	rsp := &ListResponse{}
	return rsp, t.client.Call("function", "List", request, rsp)
}

type CallRequest struct {
	// Name of the function
	Name string `json:"name"`
	// Request body that will be passed to the function
	Request map[string]interface{} `json:"request"`
}

type CallResponse struct {
	// Response body that the function returned
	Response map[string]interface{} `json:"response"`
}

type DeleteRequest struct {
	Name    string `json:"name"`
	Project string `json:"project"`
}

type DeleteResponse struct {
}

type DeployRequest struct {
	// entry point, ie. handler name in the source code
	// if not provided, defaults to the name parameter
	Entrypoint string `json:"entrypoint"`
	// function name
	Name string `json:"name"`
	// project is used for namespacing your functions
	// optional. defaults to "default".
	Project string `json:"project"`
	// github url to repo
	Repo string `json:"repo"`
	// optional subfolder path
	Subfolder string `json:"subfolder"`
}

type DeployResponse struct {
}

type Func struct {
	// name of handler in source code
	Entrypoint string `json:"entrypoint"`
	// function name
	Name string `json:"name"`
	// project of function, optional
	// defaults to literal "default"
	// used to namespace functions
	Project string `json:"project"`
	// git repo address
	Repo string `json:"repo"`
	// subfolder path to entrypoint
	Subfolder string `json:"subfolder"`
}

type ListRequest struct {
	// optional
	Project string `json:"project"`
}

type ListResponse struct {
	Functions []Func `json:"functions"`
}
