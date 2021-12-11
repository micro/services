// Code generated by protoc-gen-micro. DO NOT EDIT.
// source: proto/avatar.proto

package avatar

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

// Api Endpoints for Avatar service

func NewAvatarEndpoints() []*api.Endpoint {
	return []*api.Endpoint{}
}

// Client API for Avatar service

type AvatarService interface {
	Generate(ctx context.Context, in *GenerateRequest, opts ...client.CallOption) (*GenerateResponse, error)
}

type avatarService struct {
	c    client.Client
	name string
}

func NewAvatarService(name string, c client.Client) AvatarService {
	return &avatarService{
		c:    c,
		name: name,
	}
}

func (c *avatarService) Generate(ctx context.Context, in *GenerateRequest, opts ...client.CallOption) (*GenerateResponse, error) {
	req := c.c.NewRequest(c.name, "Avatar.Generate", in)
	out := new(GenerateResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for Avatar service

type AvatarHandler interface {
	Generate(context.Context, *GenerateRequest, *GenerateResponse) error
}

func RegisterAvatarHandler(s server.Server, hdlr AvatarHandler, opts ...server.HandlerOption) error {
	type avatar interface {
		Generate(ctx context.Context, in *GenerateRequest, out *GenerateResponse) error
	}
	type Avatar struct {
		avatar
	}
	h := &avatarHandler{hdlr}
	return s.Handle(s.NewHandler(&Avatar{h}, opts...))
}

type avatarHandler struct {
	AvatarHandler
}

func (h *avatarHandler) Generate(ctx context.Context, in *GenerateRequest, out *GenerateResponse) error {
	return h.AvatarHandler.Generate(ctx, in, out)
}
