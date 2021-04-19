// Code generated by protoc-gen-micro. DO NOT EDIT.
// source: proto/streams.proto

package streams

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

// Api Endpoints for Streams service

func NewStreamsEndpoints() []*api.Endpoint {
	return []*api.Endpoint{}
}

// Client API for Streams service

type StreamsService interface {
	Publish(ctx context.Context, in *Message, opts ...client.CallOption) (*PublishResponse, error)
	Subscribe(ctx context.Context, in *SubscribeRequest, opts ...client.CallOption) (Streams_SubscribeService, error)
	Token(ctx context.Context, in *TokenRequest, opts ...client.CallOption) (*TokenResponse, error)
}

type streamsService struct {
	c    client.Client
	name string
}

func NewStreamsService(name string, c client.Client) StreamsService {
	return &streamsService{
		c:    c,
		name: name,
	}
}

func (c *streamsService) Publish(ctx context.Context, in *Message, opts ...client.CallOption) (*PublishResponse, error) {
	req := c.c.NewRequest(c.name, "Streams.Publish", in)
	out := new(PublishResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *streamsService) Subscribe(ctx context.Context, in *SubscribeRequest, opts ...client.CallOption) (Streams_SubscribeService, error) {
	req := c.c.NewRequest(c.name, "Streams.Subscribe", &SubscribeRequest{})
	stream, err := c.c.Stream(ctx, req, opts...)
	if err != nil {
		return nil, err
	}
	if err := stream.Send(in); err != nil {
		return nil, err
	}
	return &streamsServiceSubscribe{stream}, nil
}

type Streams_SubscribeService interface {
	Context() context.Context
	SendMsg(interface{}) error
	RecvMsg(interface{}) error
	Close() error
	Recv() (*Message, error)
}

type streamsServiceSubscribe struct {
	stream client.Stream
}

func (x *streamsServiceSubscribe) Close() error {
	return x.stream.Close()
}

func (x *streamsServiceSubscribe) Context() context.Context {
	return x.stream.Context()
}

func (x *streamsServiceSubscribe) SendMsg(m interface{}) error {
	return x.stream.Send(m)
}

func (x *streamsServiceSubscribe) RecvMsg(m interface{}) error {
	return x.stream.Recv(m)
}

func (x *streamsServiceSubscribe) Recv() (*Message, error) {
	m := new(Message)
	err := x.stream.Recv(m)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func (c *streamsService) Token(ctx context.Context, in *TokenRequest, opts ...client.CallOption) (*TokenResponse, error) {
	req := c.c.NewRequest(c.name, "Streams.Token", in)
	out := new(TokenResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for Streams service

type StreamsHandler interface {
	Publish(context.Context, *Message, *PublishResponse) error
	Subscribe(context.Context, *SubscribeRequest, Streams_SubscribeStream) error
	Token(context.Context, *TokenRequest, *TokenResponse) error
}

func RegisterStreamsHandler(s server.Server, hdlr StreamsHandler, opts ...server.HandlerOption) error {
	type streams interface {
		Publish(ctx context.Context, in *Message, out *PublishResponse) error
		Subscribe(ctx context.Context, stream server.Stream) error
		Token(ctx context.Context, in *TokenRequest, out *TokenResponse) error
	}
	type Streams struct {
		streams
	}
	h := &streamsHandler{hdlr}
	return s.Handle(s.NewHandler(&Streams{h}, opts...))
}

type streamsHandler struct {
	StreamsHandler
}

func (h *streamsHandler) Publish(ctx context.Context, in *Message, out *PublishResponse) error {
	return h.StreamsHandler.Publish(ctx, in, out)
}

func (h *streamsHandler) Subscribe(ctx context.Context, stream server.Stream) error {
	m := new(SubscribeRequest)
	if err := stream.Recv(m); err != nil {
		return err
	}
	return h.StreamsHandler.Subscribe(ctx, m, &streamsSubscribeStream{stream})
}

type Streams_SubscribeStream interface {
	Context() context.Context
	SendMsg(interface{}) error
	RecvMsg(interface{}) error
	Close() error
	Send(*Message) error
}

type streamsSubscribeStream struct {
	stream server.Stream
}

func (x *streamsSubscribeStream) Close() error {
	return x.stream.Close()
}

func (x *streamsSubscribeStream) Context() context.Context {
	return x.stream.Context()
}

func (x *streamsSubscribeStream) SendMsg(m interface{}) error {
	return x.stream.Send(m)
}

func (x *streamsSubscribeStream) RecvMsg(m interface{}) error {
	return x.stream.Recv(m)
}

func (x *streamsSubscribeStream) Send(m *Message) error {
	return x.stream.Send(m)
}

func (h *streamsHandler) Token(ctx context.Context, in *TokenRequest, out *TokenResponse) error {
	return h.StreamsHandler.Token(ctx, in, out)
}
