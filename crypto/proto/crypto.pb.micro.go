// Code generated by protoc-gen-micro. DO NOT EDIT.
// source: proto/crypto.proto

package crypto

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

// Api Endpoints for Crypto service

func NewCryptoEndpoints() []*api.Endpoint {
	return []*api.Endpoint{}
}

// Client API for Crypto service

type CryptoService interface {
	News(ctx context.Context, in *NewsRequest, opts ...client.CallOption) (*NewsResponse, error)
	Quote(ctx context.Context, in *QuoteRequest, opts ...client.CallOption) (*QuoteResponse, error)
	Price(ctx context.Context, in *PriceRequest, opts ...client.CallOption) (*PriceResponse, error)
	History(ctx context.Context, in *HistoryRequest, opts ...client.CallOption) (*HistoryResponse, error)
	Symbols(ctx context.Context, in *SymbolsRequest, opts ...client.CallOption) (*SymbolsResponse, error)
}

type cryptoService struct {
	c    client.Client
	name string
}

func NewCryptoService(name string, c client.Client) CryptoService {
	return &cryptoService{
		c:    c,
		name: name,
	}
}

func (c *cryptoService) News(ctx context.Context, in *NewsRequest, opts ...client.CallOption) (*NewsResponse, error) {
	req := c.c.NewRequest(c.name, "Crypto.News", in)
	out := new(NewsResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *cryptoService) Quote(ctx context.Context, in *QuoteRequest, opts ...client.CallOption) (*QuoteResponse, error) {
	req := c.c.NewRequest(c.name, "Crypto.Quote", in)
	out := new(QuoteResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *cryptoService) Price(ctx context.Context, in *PriceRequest, opts ...client.CallOption) (*PriceResponse, error) {
	req := c.c.NewRequest(c.name, "Crypto.Price", in)
	out := new(PriceResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *cryptoService) History(ctx context.Context, in *HistoryRequest, opts ...client.CallOption) (*HistoryResponse, error) {
	req := c.c.NewRequest(c.name, "Crypto.History", in)
	out := new(HistoryResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *cryptoService) Symbols(ctx context.Context, in *SymbolsRequest, opts ...client.CallOption) (*SymbolsResponse, error) {
	req := c.c.NewRequest(c.name, "Crypto.Symbols", in)
	out := new(SymbolsResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for Crypto service

type CryptoHandler interface {
	News(context.Context, *NewsRequest, *NewsResponse) error
	Quote(context.Context, *QuoteRequest, *QuoteResponse) error
	Price(context.Context, *PriceRequest, *PriceResponse) error
	History(context.Context, *HistoryRequest, *HistoryResponse) error
	Symbols(context.Context, *SymbolsRequest, *SymbolsResponse) error
}

func RegisterCryptoHandler(s server.Server, hdlr CryptoHandler, opts ...server.HandlerOption) error {
	type crypto interface {
		News(ctx context.Context, in *NewsRequest, out *NewsResponse) error
		Quote(ctx context.Context, in *QuoteRequest, out *QuoteResponse) error
		Price(ctx context.Context, in *PriceRequest, out *PriceResponse) error
		History(ctx context.Context, in *HistoryRequest, out *HistoryResponse) error
		Symbols(ctx context.Context, in *SymbolsRequest, out *SymbolsResponse) error
	}
	type Crypto struct {
		crypto
	}
	h := &cryptoHandler{hdlr}
	return s.Handle(s.NewHandler(&Crypto{h}, opts...))
}

type cryptoHandler struct {
	CryptoHandler
}

func (h *cryptoHandler) News(ctx context.Context, in *NewsRequest, out *NewsResponse) error {
	return h.CryptoHandler.News(ctx, in, out)
}

func (h *cryptoHandler) Quote(ctx context.Context, in *QuoteRequest, out *QuoteResponse) error {
	return h.CryptoHandler.Quote(ctx, in, out)
}

func (h *cryptoHandler) Price(ctx context.Context, in *PriceRequest, out *PriceResponse) error {
	return h.CryptoHandler.Price(ctx, in, out)
}

func (h *cryptoHandler) History(ctx context.Context, in *HistoryRequest, out *HistoryResponse) error {
	return h.CryptoHandler.History(ctx, in, out)
}

func (h *cryptoHandler) Symbols(ctx context.Context, in *SymbolsRequest, out *SymbolsResponse) error {
	return h.CryptoHandler.Symbols(ctx, in, out)
}
