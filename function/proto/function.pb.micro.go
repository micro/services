// Code generated by protoc-gen-micro. DO NOT EDIT.
// source: proto/function.proto

package function

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	_ "google.golang.org/protobuf/types/known/structpb"
	math "math"
)

import (
	context "context"
	api "github.com/micro/micro/v3/service/api"
	client "github.com/micro/micro/v3/service/client"
	server "github.com/micro/micro/v3/service/server"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

// Reference imports to suppress errors if they are not otherwise used.
var _ api.Endpoint
var _ context.Context
var _ client.Option
var _ server.Option

// Api Endpoints for Function service

func NewFunctionEndpoints() []*api.Endpoint {
	return []*api.Endpoint{}
}

// Client API for Function service

type FunctionService interface {
	Call(ctx context.Context, in *CallRequest, opts ...client.CallOption) (*CallResponse, error)
	Create(ctx context.Context, in *CreateRequest, opts ...client.CallOption) (*CreateResponse, error)
	Delete(ctx context.Context, in *DeleteRequest, opts ...client.CallOption) (*DeleteResponse, error)
	Describe(ctx context.Context, in *DescribeRequest, opts ...client.CallOption) (*DescribeResponse, error)
	List(ctx context.Context, in *ListRequest, opts ...client.CallOption) (*ListResponse, error)
}

type functionService struct {
	c    client.Client
	name string
}

func NewFunctionService(name string, c client.Client) FunctionService {
	return &functionService{
		c:    c,
		name: name,
	}
}

func (c *functionService) Call(ctx context.Context, in *CallRequest, opts ...client.CallOption) (*CallResponse, error) {
	req := c.c.NewRequest(c.name, "Function.Call", in)
	out := new(CallResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *functionService) Create(ctx context.Context, in *CreateRequest, opts ...client.CallOption) (*CreateResponse, error) {
	req := c.c.NewRequest(c.name, "Function.Create", in)
	out := new(CreateResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *functionService) Delete(ctx context.Context, in *DeleteRequest, opts ...client.CallOption) (*DeleteResponse, error) {
	req := c.c.NewRequest(c.name, "Function.Delete", in)
	out := new(DeleteResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *functionService) Describe(ctx context.Context, in *DescribeRequest, opts ...client.CallOption) (*DescribeResponse, error) {
	req := c.c.NewRequest(c.name, "Function.Describe", in)
	out := new(DescribeResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *functionService) List(ctx context.Context, in *ListRequest, opts ...client.CallOption) (*ListResponse, error) {
	req := c.c.NewRequest(c.name, "Function.List", in)
	out := new(ListResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for Function service

type FunctionHandler interface {
	Call(context.Context, *CallRequest, *CallResponse) error
	Create(context.Context, *CreateRequest, *CreateResponse) error
	Delete(context.Context, *DeleteRequest, *DeleteResponse) error
	Describe(context.Context, *DescribeRequest, *DescribeResponse) error
	List(context.Context, *ListRequest, *ListResponse) error
}

func RegisterFunctionHandler(s server.Server, hdlr FunctionHandler, opts ...server.HandlerOption) error {
	type function interface {
		Call(ctx context.Context, in *CallRequest, out *CallResponse) error
		Create(ctx context.Context, in *CreateRequest, out *CreateResponse) error
		Delete(ctx context.Context, in *DeleteRequest, out *DeleteResponse) error
		Describe(ctx context.Context, in *DescribeRequest, out *DescribeResponse) error
		List(ctx context.Context, in *ListRequest, out *ListResponse) error
	}
	type Function struct {
		function
	}
	h := &functionHandler{hdlr}
	return s.Handle(s.NewHandler(&Function{h}, opts...))
}

type functionHandler struct {
	FunctionHandler
}

func (h *functionHandler) Call(ctx context.Context, in *CallRequest, out *CallResponse) error {
	return h.FunctionHandler.Call(ctx, in, out)
}

func (h *functionHandler) Create(ctx context.Context, in *CreateRequest, out *CreateResponse) error {
	return h.FunctionHandler.Create(ctx, in, out)
}

func (h *functionHandler) Delete(ctx context.Context, in *DeleteRequest, out *DeleteResponse) error {
	return h.FunctionHandler.Delete(ctx, in, out)
}

func (h *functionHandler) Describe(ctx context.Context, in *DescribeRequest, out *DescribeResponse) error {
	return h.FunctionHandler.Describe(ctx, in, out)
}

func (h *functionHandler) List(ctx context.Context, in *ListRequest, out *ListResponse) error {
	return h.FunctionHandler.List(ctx, in, out)
}
