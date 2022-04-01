// Code generated by protoc-gen-micro. DO NOT EDIT.
// source: proto/price.proto

package price

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

// Api Endpoints for Price service

func NewPriceEndpoints() []*api.Endpoint {
	return []*api.Endpoint{}
}

// Client API for Price service

type PriceService interface {
	Add(ctx context.Context, in *AddRequest, opts ...client.CallOption) (*AddResponse, error)
	Get(ctx context.Context, in *GetRequest, opts ...client.CallOption) (*GetResponse, error)
	List(ctx context.Context, in *ListRequest, opts ...client.CallOption) (*ListResponse, error)
	Index(ctx context.Context, in *IndexRequest, opts ...client.CallOption) (*IndexResponse, error)
	Report(ctx context.Context, in *ReportRequest, opts ...client.CallOption) (*ReportResponse, error)
}

type priceService struct {
	c    client.Client
	name string
}

func NewPriceService(name string, c client.Client) PriceService {
	return &priceService{
		c:    c,
		name: name,
	}
}

func (c *priceService) Add(ctx context.Context, in *AddRequest, opts ...client.CallOption) (*AddResponse, error) {
	req := c.c.NewRequest(c.name, "Price.Add", in)
	out := new(AddResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *priceService) Get(ctx context.Context, in *GetRequest, opts ...client.CallOption) (*GetResponse, error) {
	req := c.c.NewRequest(c.name, "Price.Get", in)
	out := new(GetResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *priceService) List(ctx context.Context, in *ListRequest, opts ...client.CallOption) (*ListResponse, error) {
	req := c.c.NewRequest(c.name, "Price.List", in)
	out := new(ListResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *priceService) Index(ctx context.Context, in *IndexRequest, opts ...client.CallOption) (*IndexResponse, error) {
	req := c.c.NewRequest(c.name, "Price.Index", in)
	out := new(IndexResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *priceService) Report(ctx context.Context, in *ReportRequest, opts ...client.CallOption) (*ReportResponse, error) {
	req := c.c.NewRequest(c.name, "Price.Report", in)
	out := new(ReportResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for Price service

type PriceHandler interface {
	Add(context.Context, *AddRequest, *AddResponse) error
	Get(context.Context, *GetRequest, *GetResponse) error
	List(context.Context, *ListRequest, *ListResponse) error
	Index(context.Context, *IndexRequest, *IndexResponse) error
	Report(context.Context, *ReportRequest, *ReportResponse) error
}

func RegisterPriceHandler(s server.Server, hdlr PriceHandler, opts ...server.HandlerOption) error {
	type price interface {
		Add(ctx context.Context, in *AddRequest, out *AddResponse) error
		Get(ctx context.Context, in *GetRequest, out *GetResponse) error
		List(ctx context.Context, in *ListRequest, out *ListResponse) error
		Index(ctx context.Context, in *IndexRequest, out *IndexResponse) error
		Report(ctx context.Context, in *ReportRequest, out *ReportResponse) error
	}
	type Price struct {
		price
	}
	h := &priceHandler{hdlr}
	return s.Handle(s.NewHandler(&Price{h}, opts...))
}

type priceHandler struct {
	PriceHandler
}

func (h *priceHandler) Add(ctx context.Context, in *AddRequest, out *AddResponse) error {
	return h.PriceHandler.Add(ctx, in, out)
}

func (h *priceHandler) Get(ctx context.Context, in *GetRequest, out *GetResponse) error {
	return h.PriceHandler.Get(ctx, in, out)
}

func (h *priceHandler) List(ctx context.Context, in *ListRequest, out *ListResponse) error {
	return h.PriceHandler.List(ctx, in, out)
}

func (h *priceHandler) Index(ctx context.Context, in *IndexRequest, out *IndexResponse) error {
	return h.PriceHandler.Index(ctx, in, out)
}

func (h *priceHandler) Report(ctx context.Context, in *ReportRequest, out *ReportResponse) error {
	return h.PriceHandler.Report(ctx, in, out)
}