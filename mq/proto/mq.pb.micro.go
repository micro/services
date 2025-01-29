// Code generated by protoc-gen-micro. DO NOT EDIT.
// source: proto/mq.proto

package mq

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	_ "github.com/golang/protobuf/ptypes/struct"
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

// Client API for Mq service

type MqService interface {
	Publish(ctx context.Context, in *PublishRequest, opts ...client.CallOption) (*PublishResponse, error)
	Subscribe(ctx context.Context, in *SubscribeRequest, opts ...client.CallOption) (Mq_SubscribeService, error)
}

type mqService struct {
	c    client.Client
	name string
}

func NewMqService(name string, c client.Client) MqService {
	return &mqService{
		c:    c,
		name: name,
	}
}

func (c *mqService) Publish(ctx context.Context, in *PublishRequest, opts ...client.CallOption) (*PublishResponse, error) {
	req := c.c.NewRequest(c.name, "Mq.Publish", in)
	out := new(PublishResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *mqService) Subscribe(ctx context.Context, in *SubscribeRequest, opts ...client.CallOption) (Mq_SubscribeService, error) {
	req := c.c.NewRequest(c.name, "Mq.Subscribe", &SubscribeRequest{})
	stream, err := c.c.Stream(ctx, req, opts...)
	if err != nil {
		return nil, err
	}
	if err := stream.Send(in); err != nil {
		return nil, err
	}
	return &mqServiceSubscribe{stream}, nil
}

type Mq_SubscribeService interface {
	Context() context.Context
	SendMsg(interface{}) error
	RecvMsg(interface{}) error
	Close() error
	Recv() (*SubscribeResponse, error)
}

type mqServiceSubscribe struct {
	stream client.Stream
}

func (x *mqServiceSubscribe) Close() error {
	return x.stream.Close()
}

func (x *mqServiceSubscribe) Context() context.Context {
	return x.stream.Context()
}

func (x *mqServiceSubscribe) SendMsg(m interface{}) error {
	return x.stream.Send(m)
}

func (x *mqServiceSubscribe) RecvMsg(m interface{}) error {
	return x.stream.Recv(m)
}

func (x *mqServiceSubscribe) Recv() (*SubscribeResponse, error) {
	m := new(SubscribeResponse)
	err := x.stream.Recv(m)
	if err != nil {
		return nil, err
	}
	return m, nil
}

// Server API for Mq service

type MqHandler interface {
	Publish(context.Context, *PublishRequest, *PublishResponse) error
	Subscribe(context.Context, *SubscribeRequest, Mq_SubscribeStream) error
}

func RegisterMqHandler(s server.Server, hdlr MqHandler, opts ...server.HandlerOption) error {
	type mq interface {
		Publish(ctx context.Context, in *PublishRequest, out *PublishResponse) error
		Subscribe(ctx context.Context, stream server.Stream) error
	}
	type Mq struct {
		mq
	}
	h := &mqHandler{hdlr}
	return s.Handle(s.NewHandler(&Mq{h}, opts...))
}

type mqHandler struct {
	MqHandler
}

func (h *mqHandler) Publish(ctx context.Context, in *PublishRequest, out *PublishResponse) error {
	return h.MqHandler.Publish(ctx, in, out)
}

func (h *mqHandler) Subscribe(ctx context.Context, stream server.Stream) error {
	m := new(SubscribeRequest)
	if err := stream.Recv(m); err != nil {
		return err
	}
	return h.MqHandler.Subscribe(ctx, m, &mqSubscribeStream{stream})
}

type Mq_SubscribeStream interface {
	Context() context.Context
	SendMsg(interface{}) error
	RecvMsg(interface{}) error
	Close() error
	Send(*SubscribeResponse) error
}

type mqSubscribeStream struct {
	stream server.Stream
}

func (x *mqSubscribeStream) Close() error {
	return x.stream.Close()
}

func (x *mqSubscribeStream) Context() context.Context {
	return x.stream.Context()
}

func (x *mqSubscribeStream) SendMsg(m interface{}) error {
	return x.stream.Send(m)
}

func (x *mqSubscribeStream) RecvMsg(m interface{}) error {
	return x.stream.Recv(m)
}

func (x *mqSubscribeStream) Send(m *SubscribeResponse) error {
	return x.stream.Send(m)
}
