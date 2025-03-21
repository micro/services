// Code generated by protoc-gen-micro. DO NOT EDIT.
// source: proto/youtube.proto

package youtube

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

// Client API for Youtube service

type YoutubeService interface {
	Search(ctx context.Context, in *SearchRequest, opts ...client.CallOption) (*SearchResponse, error)
	Embed(ctx context.Context, in *EmbedRequest, opts ...client.CallOption) (*EmbedResponse, error)
}

type youtubeService struct {
	c    client.Client
	name string
}

func NewYoutubeService(name string, c client.Client) YoutubeService {
	return &youtubeService{
		c:    c,
		name: name,
	}
}

func (c *youtubeService) Search(ctx context.Context, in *SearchRequest, opts ...client.CallOption) (*SearchResponse, error) {
	req := c.c.NewRequest(c.name, "Youtube.Search", in)
	out := new(SearchResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *youtubeService) Embed(ctx context.Context, in *EmbedRequest, opts ...client.CallOption) (*EmbedResponse, error) {
	req := c.c.NewRequest(c.name, "Youtube.Embed", in)
	out := new(EmbedResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for Youtube service

type YoutubeHandler interface {
	Search(context.Context, *SearchRequest, *SearchResponse) error
	Embed(context.Context, *EmbedRequest, *EmbedResponse) error
}

func RegisterYoutubeHandler(s server.Server, hdlr YoutubeHandler, opts ...server.HandlerOption) error {
	type youtube interface {
		Search(ctx context.Context, in *SearchRequest, out *SearchResponse) error
		Embed(ctx context.Context, in *EmbedRequest, out *EmbedResponse) error
	}
	type Youtube struct {
		youtube
	}
	h := &youtubeHandler{hdlr}
	return s.Handle(s.NewHandler(&Youtube{h}, opts...))
}

type youtubeHandler struct {
	YoutubeHandler
}

func (h *youtubeHandler) Search(ctx context.Context, in *SearchRequest, out *SearchResponse) error {
	return h.YoutubeHandler.Search(ctx, in, out)
}

func (h *youtubeHandler) Embed(ctx context.Context, in *EmbedRequest, out *EmbedResponse) error {
	return h.YoutubeHandler.Embed(ctx, in, out)
}
