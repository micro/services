// Code generated by protoc-gen-micro. DO NOT EDIT.
// source: proto/google.proto

package google

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	math "math"
)

import (
	context "context"
	client "github.com/micro/micro/v5/service/client"
	server "github.com/micro/micro/v5/service/server"
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
var _ context.Context
var _ client.Option
var _ server.Option

// Client API for Google service

type GoogleService interface {
	Search(ctx context.Context, in *SearchRequest, opts ...client.CallOption) (*SearchResponse, error)
}

type googleService struct {
	c    client.Client
	name string
}

func NewGoogleService(name string, c client.Client) GoogleService {
	return &googleService{
		c:    c,
		name: name,
	}
}

func (c *googleService) Search(ctx context.Context, in *SearchRequest, opts ...client.CallOption) (*SearchResponse, error) {
	req := c.c.NewRequest(c.name, "Google.Search", in)
	out := new(SearchResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for Google service

type GoogleHandler interface {
	Search(context.Context, *SearchRequest, *SearchResponse) error
}

func RegisterGoogleHandler(s server.Server, hdlr GoogleHandler, opts ...server.HandlerOption) error {
	type google interface {
		Search(ctx context.Context, in *SearchRequest, out *SearchResponse) error
	}
	type Google struct {
		google
	}
	h := &googleHandler{hdlr}
	return s.Handle(s.NewHandler(&Google{h}, opts...))
}

type googleHandler struct {
	GoogleHandler
}

func (h *googleHandler) Search(ctx context.Context, in *SearchRequest, out *SearchResponse) error {
	return h.GoogleHandler.Search(ctx, in, out)
}
