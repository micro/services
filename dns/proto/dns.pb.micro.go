// Code generated by protoc-gen-micro. DO NOT EDIT.
// source: proto/dns.proto

package dns

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

// Client API for Dns service

type DnsService interface {
	Query(ctx context.Context, in *QueryRequest, opts ...client.CallOption) (*QueryResponse, error)
	Whois(ctx context.Context, in *WhoisRequest, opts ...client.CallOption) (*WhoisResponse, error)
}

type dnsService struct {
	c    client.Client
	name string
}

func NewDnsService(name string, c client.Client) DnsService {
	return &dnsService{
		c:    c,
		name: name,
	}
}

func (c *dnsService) Query(ctx context.Context, in *QueryRequest, opts ...client.CallOption) (*QueryResponse, error) {
	req := c.c.NewRequest(c.name, "Dns.Query", in)
	out := new(QueryResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *dnsService) Whois(ctx context.Context, in *WhoisRequest, opts ...client.CallOption) (*WhoisResponse, error) {
	req := c.c.NewRequest(c.name, "Dns.Whois", in)
	out := new(WhoisResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for Dns service

type DnsHandler interface {
	Query(context.Context, *QueryRequest, *QueryResponse) error
	Whois(context.Context, *WhoisRequest, *WhoisResponse) error
}

func RegisterDnsHandler(s server.Server, hdlr DnsHandler, opts ...server.HandlerOption) error {
	type dns interface {
		Query(ctx context.Context, in *QueryRequest, out *QueryResponse) error
		Whois(ctx context.Context, in *WhoisRequest, out *WhoisResponse) error
	}
	type Dns struct {
		dns
	}
	h := &dnsHandler{hdlr}
	return s.Handle(s.NewHandler(&Dns{h}, opts...))
}

type dnsHandler struct {
	DnsHandler
}

func (h *dnsHandler) Query(ctx context.Context, in *QueryRequest, out *QueryResponse) error {
	return h.DnsHandler.Query(ctx, in, out)
}

func (h *dnsHandler) Whois(ctx context.Context, in *WhoisRequest, out *WhoisResponse) error {
	return h.DnsHandler.Whois(ctx, in, out)
}
