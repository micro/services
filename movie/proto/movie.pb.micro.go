// Code generated by protoc-gen-micro. DO NOT EDIT.
// source: proto/movie.proto

package movie

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

// Client API for Movie service

type MovieService interface {
	Search(ctx context.Context, in *SearchRequest, opts ...client.CallOption) (*SearchResponse, error)
}

type movieService struct {
	c    client.Client
	name string
}

func NewMovieService(name string, c client.Client) MovieService {
	return &movieService{
		c:    c,
		name: name,
	}
}

func (c *movieService) Search(ctx context.Context, in *SearchRequest, opts ...client.CallOption) (*SearchResponse, error) {
	req := c.c.NewRequest(c.name, "Movie.Search", in)
	out := new(SearchResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for Movie service

type MovieHandler interface {
	Search(context.Context, *SearchRequest, *SearchResponse) error
}

func RegisterMovieHandler(s server.Server, hdlr MovieHandler, opts ...server.HandlerOption) error {
	type movie interface {
		Search(ctx context.Context, in *SearchRequest, out *SearchResponse) error
	}
	type Movie struct {
		movie
	}
	h := &movieHandler{hdlr}
	return s.Handle(s.NewHandler(&Movie{h}, opts...))
}

type movieHandler struct {
	MovieHandler
}

func (h *movieHandler) Search(ctx context.Context, in *SearchRequest, out *SearchResponse) error {
	return h.MovieHandler.Search(ctx, in, out)
}
