// Code generated by protoc-gen-micro. DO NOT EDIT.
// source: proto/twitter.proto

package twitter

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

// Api Endpoints for Twitter service

func NewTwitterEndpoints() []*api.Endpoint {
	return []*api.Endpoint{}
}

// Client API for Twitter service

type TwitterService interface {
	Timeline(ctx context.Context, in *TimelineRequest, opts ...client.CallOption) (*TimelineResponse, error)
	Search(ctx context.Context, in *SearchRequest, opts ...client.CallOption) (*SearchResponse, error)
	User(ctx context.Context, in *UserRequest, opts ...client.CallOption) (*UserResponse, error)
}

type twitterService struct {
	c    client.Client
	name string
}

func NewTwitterService(name string, c client.Client) TwitterService {
	return &twitterService{
		c:    c,
		name: name,
	}
}

func (c *twitterService) Timeline(ctx context.Context, in *TimelineRequest, opts ...client.CallOption) (*TimelineResponse, error) {
	req := c.c.NewRequest(c.name, "Twitter.Timeline", in)
	out := new(TimelineResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *twitterService) Search(ctx context.Context, in *SearchRequest, opts ...client.CallOption) (*SearchResponse, error) {
	req := c.c.NewRequest(c.name, "Twitter.Search", in)
	out := new(SearchResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *twitterService) User(ctx context.Context, in *UserRequest, opts ...client.CallOption) (*UserResponse, error) {
	req := c.c.NewRequest(c.name, "Twitter.User", in)
	out := new(UserResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for Twitter service

type TwitterHandler interface {
	Timeline(context.Context, *TimelineRequest, *TimelineResponse) error
	Search(context.Context, *SearchRequest, *SearchResponse) error
	User(context.Context, *UserRequest, *UserResponse) error
}

func RegisterTwitterHandler(s server.Server, hdlr TwitterHandler, opts ...server.HandlerOption) error {
	type twitter interface {
		Timeline(ctx context.Context, in *TimelineRequest, out *TimelineResponse) error
		Search(ctx context.Context, in *SearchRequest, out *SearchResponse) error
		User(ctx context.Context, in *UserRequest, out *UserResponse) error
	}
	type Twitter struct {
		twitter
	}
	h := &twitterHandler{hdlr}
	return s.Handle(s.NewHandler(&Twitter{h}, opts...))
}

type twitterHandler struct {
	TwitterHandler
}

func (h *twitterHandler) Timeline(ctx context.Context, in *TimelineRequest, out *TimelineResponse) error {
	return h.TwitterHandler.Timeline(ctx, in, out)
}

func (h *twitterHandler) Search(ctx context.Context, in *SearchRequest, out *SearchResponse) error {
	return h.TwitterHandler.Search(ctx, in, out)
}

func (h *twitterHandler) User(ctx context.Context, in *UserRequest, out *UserResponse) error {
	return h.TwitterHandler.User(ctx, in, out)
}
