// Code generated by protoc-gen-micro. DO NOT EDIT.
// source: proto/password.proto

package password

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

// Client API for Password service

type PasswordService interface {
	Generate(ctx context.Context, in *GenerateRequest, opts ...client.CallOption) (*GenerateResponse, error)
}

type passwordService struct {
	c    client.Client
	name string
}

func NewPasswordService(name string, c client.Client) PasswordService {
	return &passwordService{
		c:    c,
		name: name,
	}
}

func (c *passwordService) Generate(ctx context.Context, in *GenerateRequest, opts ...client.CallOption) (*GenerateResponse, error) {
	req := c.c.NewRequest(c.name, "Password.Generate", in)
	out := new(GenerateResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for Password service

type PasswordHandler interface {
	Generate(context.Context, *GenerateRequest, *GenerateResponse) error
}

func RegisterPasswordHandler(s server.Server, hdlr PasswordHandler, opts ...server.HandlerOption) error {
	type password interface {
		Generate(ctx context.Context, in *GenerateRequest, out *GenerateResponse) error
	}
	type Password struct {
		password
	}
	h := &passwordHandler{hdlr}
	return s.Handle(s.NewHandler(&Password{h}, opts...))
}

type passwordHandler struct {
	PasswordHandler
}

func (h *passwordHandler) Generate(ctx context.Context, in *GenerateRequest, out *GenerateResponse) error {
	return h.PasswordHandler.Generate(ctx, in, out)
}
