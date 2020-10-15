// Code generated by protoc-gen-micro. DO NOT EDIT.
// source: proto/test.proto

package test

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

// Api Endpoints for Test service

func NewTestEndpoints() []*api.Endpoint {
	return []*api.Endpoint{}
}

// Client API for Test service

type TestService interface {
	// Call handles a single request
	Call(ctx context.Context, in *Request, opts ...client.CallOption) (*Response, error)
	// Stream handles a streaming request
	Stream(ctx context.Context, in *Request, opts ...client.CallOption) (Test_StreamService, error)
	// Config tests the usage of config service
	Config(ctx context.Context, in *Request, opts ...client.CallOption) (*Response, error)
	// Store tests the usage of the store service
	Store(ctx context.Context, in *Request, opts ...client.CallOption) (*Response, error)
	// Events tests the usage of the events service
	Events(ctx context.Context, in *Request, opts ...client.CallOption) (*Response, error)
	// Broker tests the usage of the broker service
	Broker(ctx context.Context, in *Request, opts ...client.CallOption) (*Response, error)
	// BlobStore tests the usage of the blob store
	BlobStore(ctx context.Context, in *Request, opts ...client.CallOption) (*Response, error)
	// Logger tests the usage of the service logger
	Logger(ctx context.Context, in *Request, opts ...client.CallOption) (*Response, error)
}

type testService struct {
	c    client.Client
	name string
}

func NewTestService(name string, c client.Client) TestService {
	return &testService{
		c:    c,
		name: name,
	}
}

func (c *testService) Call(ctx context.Context, in *Request, opts ...client.CallOption) (*Response, error) {
	req := c.c.NewRequest(c.name, "Test.Call", in)
	out := new(Response)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *testService) Stream(ctx context.Context, in *Request, opts ...client.CallOption) (Test_StreamService, error) {
	req := c.c.NewRequest(c.name, "Test.Stream", &Request{})
	stream, err := c.c.Stream(ctx, req, opts...)
	if err != nil {
		return nil, err
	}
	if err := stream.Send(in); err != nil {
		return nil, err
	}
	return &testServiceStream{stream}, nil
}

type Test_StreamService interface {
	Context() context.Context
	SendMsg(interface{}) error
	RecvMsg(interface{}) error
	Close() error
	Recv() (*Response, error)
}

type testServiceStream struct {
	stream client.Stream
}

func (x *testServiceStream) Close() error {
	return x.stream.Close()
}

func (x *testServiceStream) Context() context.Context {
	return x.stream.Context()
}

func (x *testServiceStream) SendMsg(m interface{}) error {
	return x.stream.Send(m)
}

func (x *testServiceStream) RecvMsg(m interface{}) error {
	return x.stream.Recv(m)
}

func (x *testServiceStream) Recv() (*Response, error) {
	m := new(Response)
	err := x.stream.Recv(m)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func (c *testService) Config(ctx context.Context, in *Request, opts ...client.CallOption) (*Response, error) {
	req := c.c.NewRequest(c.name, "Test.Config", in)
	out := new(Response)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *testService) Store(ctx context.Context, in *Request, opts ...client.CallOption) (*Response, error) {
	req := c.c.NewRequest(c.name, "Test.Store", in)
	out := new(Response)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *testService) Events(ctx context.Context, in *Request, opts ...client.CallOption) (*Response, error) {
	req := c.c.NewRequest(c.name, "Test.Events", in)
	out := new(Response)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *testService) Broker(ctx context.Context, in *Request, opts ...client.CallOption) (*Response, error) {
	req := c.c.NewRequest(c.name, "Test.Broker", in)
	out := new(Response)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *testService) BlobStore(ctx context.Context, in *Request, opts ...client.CallOption) (*Response, error) {
	req := c.c.NewRequest(c.name, "Test.BlobStore", in)
	out := new(Response)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *testService) Logger(ctx context.Context, in *Request, opts ...client.CallOption) (*Response, error) {
	req := c.c.NewRequest(c.name, "Test.Logger", in)
	out := new(Response)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for Test service

type TestHandler interface {
	// Call handles a single request
	Call(context.Context, *Request, *Response) error
	// Stream handles a streaming request
	Stream(context.Context, *Request, Test_StreamStream) error
	// Config tests the usage of config service
	Config(context.Context, *Request, *Response) error
	// Store tests the usage of the store service
	Store(context.Context, *Request, *Response) error
	// Events tests the usage of the events service
	Events(context.Context, *Request, *Response) error
	// Broker tests the usage of the broker service
	Broker(context.Context, *Request, *Response) error
	// BlobStore tests the usage of the blob store
	BlobStore(context.Context, *Request, *Response) error
	// Logger tests the usage of the service logger
	Logger(context.Context, *Request, *Response) error
}

func RegisterTestHandler(s server.Server, hdlr TestHandler, opts ...server.HandlerOption) error {
	type test interface {
		Call(ctx context.Context, in *Request, out *Response) error
		Stream(ctx context.Context, stream server.Stream) error
		Config(ctx context.Context, in *Request, out *Response) error
		Store(ctx context.Context, in *Request, out *Response) error
		Events(ctx context.Context, in *Request, out *Response) error
		Broker(ctx context.Context, in *Request, out *Response) error
		BlobStore(ctx context.Context, in *Request, out *Response) error
		Logger(ctx context.Context, in *Request, out *Response) error
	}
	type Test struct {
		test
	}
	h := &testHandler{hdlr}
	return s.Handle(s.NewHandler(&Test{h}, opts...))
}

type testHandler struct {
	TestHandler
}

func (h *testHandler) Call(ctx context.Context, in *Request, out *Response) error {
	return h.TestHandler.Call(ctx, in, out)
}

func (h *testHandler) Stream(ctx context.Context, stream server.Stream) error {
	m := new(Request)
	if err := stream.Recv(m); err != nil {
		return err
	}
	return h.TestHandler.Stream(ctx, m, &testStreamStream{stream})
}

type Test_StreamStream interface {
	Context() context.Context
	SendMsg(interface{}) error
	RecvMsg(interface{}) error
	Close() error
	Send(*Response) error
}

type testStreamStream struct {
	stream server.Stream
}

func (x *testStreamStream) Close() error {
	return x.stream.Close()
}

func (x *testStreamStream) Context() context.Context {
	return x.stream.Context()
}

func (x *testStreamStream) SendMsg(m interface{}) error {
	return x.stream.Send(m)
}

func (x *testStreamStream) RecvMsg(m interface{}) error {
	return x.stream.Recv(m)
}

func (x *testStreamStream) Send(m *Response) error {
	return x.stream.Send(m)
}

func (h *testHandler) Config(ctx context.Context, in *Request, out *Response) error {
	return h.TestHandler.Config(ctx, in, out)
}

func (h *testHandler) Store(ctx context.Context, in *Request, out *Response) error {
	return h.TestHandler.Store(ctx, in, out)
}

func (h *testHandler) Events(ctx context.Context, in *Request, out *Response) error {
	return h.TestHandler.Events(ctx, in, out)
}

func (h *testHandler) Broker(ctx context.Context, in *Request, out *Response) error {
	return h.TestHandler.Broker(ctx, in, out)
}

func (h *testHandler) BlobStore(ctx context.Context, in *Request, out *Response) error {
	return h.TestHandler.BlobStore(ctx, in, out)
}

func (h *testHandler) Logger(ctx context.Context, in *Request, out *Response) error {
	return h.TestHandler.Logger(ctx, in, out)
}
