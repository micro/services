// Code generated by protoc-gen-micro. DO NOT EDIT.
// source: proto/qr.proto

package qr

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

// Client API for Qr service

type QrService interface {
	// Generate a QR code
	Generate(ctx context.Context, in *GenerateRequest, opts ...client.CallOption) (*GenerateResponse, error)
	// List your QR codes
	Codes(ctx context.Context, in *CodesRequest, opts ...client.CallOption) (*CodesResponse, error)
}

type qrService struct {
	c    client.Client
	name string
}

func NewQrService(name string, c client.Client) QrService {
	return &qrService{
		c:    c,
		name: name,
	}
}

func (c *qrService) Generate(ctx context.Context, in *GenerateRequest, opts ...client.CallOption) (*GenerateResponse, error) {
	req := c.c.NewRequest(c.name, "Qr.Generate", in)
	out := new(GenerateResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *qrService) Codes(ctx context.Context, in *CodesRequest, opts ...client.CallOption) (*CodesResponse, error) {
	req := c.c.NewRequest(c.name, "Qr.Codes", in)
	out := new(CodesResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for Qr service

type QrHandler interface {
	// Generate a QR code
	Generate(context.Context, *GenerateRequest, *GenerateResponse) error
	// List your QR codes
	Codes(context.Context, *CodesRequest, *CodesResponse) error
}

func RegisterQrHandler(s server.Server, hdlr QrHandler, opts ...server.HandlerOption) error {
	type qr interface {
		Generate(ctx context.Context, in *GenerateRequest, out *GenerateResponse) error
		Codes(ctx context.Context, in *CodesRequest, out *CodesResponse) error
	}
	type Qr struct {
		qr
	}
	h := &qrHandler{hdlr}
	return s.Handle(s.NewHandler(&Qr{h}, opts...))
}

type qrHandler struct {
	QrHandler
}

func (h *qrHandler) Generate(ctx context.Context, in *GenerateRequest, out *GenerateResponse) error {
	return h.QrHandler.Generate(ctx, in, out)
}

func (h *qrHandler) Codes(ctx context.Context, in *CodesRequest, out *CodesResponse) error {
	return h.QrHandler.Codes(ctx, in, out)
}
