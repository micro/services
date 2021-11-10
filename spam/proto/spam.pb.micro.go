// Code generated by protoc-gen-micro. DO NOT EDIT.
// source: proto/spam.proto

package spam

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
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

// Api Endpoints for Spam service

func NewSpamEndpoints() []*api.Endpoint {
	return []*api.Endpoint{}
}

// Client API for Spam service

type SpamService interface {
	Check(ctx context.Context, in *CheckRequest, opts ...client.CallOption) (*CheckResponse, error)
}

type spamService struct {
	c    client.Client
	name string
}

func NewSpamService(name string, c client.Client) SpamService {
	return &spamService{
		c:    c,
		name: name,
	}
}

func (c *spamService) Check(ctx context.Context, in *CheckRequest, opts ...client.CallOption) (*CheckResponse, error) {
	req := c.c.NewRequest(c.name, "Spam.Check", in)
	out := new(CheckResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for Spam service

type SpamHandler interface {
	Check(context.Context, *CheckRequest, *CheckResponse) error
}

func RegisterSpamHandler(s server.Server, hdlr SpamHandler, opts ...server.HandlerOption) error {
	type spam interface {
		Check(ctx context.Context, in *CheckRequest, out *CheckResponse) error
	}
	type Spam struct {
		spam
	}
	h := &spamHandler{hdlr}
	return s.Handle(s.NewHandler(&Spam{h}, opts...))
}

type spamHandler struct {
	SpamHandler
}

func (h *spamHandler) Check(ctx context.Context, in *CheckRequest, out *CheckResponse) error {
	return h.SpamHandler.Check(ctx, in, out)
}
