// Code generated by protoc-gen-micro. DO NOT EDIT.
// source: proto/vehicle.proto

package vehicle

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

// Client API for Vehicle service

type VehicleService interface {
	Lookup(ctx context.Context, in *LookupRequest, opts ...client.CallOption) (*LookupResponse, error)
}

type vehicleService struct {
	c    client.Client
	name string
}

func NewVehicleService(name string, c client.Client) VehicleService {
	return &vehicleService{
		c:    c,
		name: name,
	}
}

func (c *vehicleService) Lookup(ctx context.Context, in *LookupRequest, opts ...client.CallOption) (*LookupResponse, error) {
	req := c.c.NewRequest(c.name, "Vehicle.Lookup", in)
	out := new(LookupResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for Vehicle service

type VehicleHandler interface {
	Lookup(context.Context, *LookupRequest, *LookupResponse) error
}

func RegisterVehicleHandler(s server.Server, hdlr VehicleHandler, opts ...server.HandlerOption) error {
	type vehicle interface {
		Lookup(ctx context.Context, in *LookupRequest, out *LookupResponse) error
	}
	type Vehicle struct {
		vehicle
	}
	h := &vehicleHandler{hdlr}
	return s.Handle(s.NewHandler(&Vehicle{h}, opts...))
}

type vehicleHandler struct {
	VehicleHandler
}

func (h *vehicleHandler) Lookup(ctx context.Context, in *LookupRequest, out *LookupResponse) error {
	return h.VehicleHandler.Lookup(ctx, in, out)
}
