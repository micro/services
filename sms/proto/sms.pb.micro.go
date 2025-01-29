// Code generated by protoc-gen-micro. DO NOT EDIT.
// source: proto/sms.proto

package sms

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

// Client API for Sms service

type SmsService interface {
	Send(ctx context.Context, in *SendRequest, opts ...client.CallOption) (*SendResponse, error)
}

type smsService struct {
	c    client.Client
	name string
}

func NewSmsService(name string, c client.Client) SmsService {
	return &smsService{
		c:    c,
		name: name,
	}
}

func (c *smsService) Send(ctx context.Context, in *SendRequest, opts ...client.CallOption) (*SendResponse, error) {
	req := c.c.NewRequest(c.name, "Sms.Send", in)
	out := new(SendResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for Sms service

type SmsHandler interface {
	Send(context.Context, *SendRequest, *SendResponse) error
}

func RegisterSmsHandler(s server.Server, hdlr SmsHandler, opts ...server.HandlerOption) error {
	type sms interface {
		Send(ctx context.Context, in *SendRequest, out *SendResponse) error
	}
	type Sms struct {
		sms
	}
	h := &smsHandler{hdlr}
	return s.Handle(s.NewHandler(&Sms{h}, opts...))
}

type smsHandler struct {
	SmsHandler
}

func (h *smsHandler) Send(ctx context.Context, in *SendRequest, out *SendResponse) error {
	return h.SmsHandler.Send(ctx, in, out)
}
