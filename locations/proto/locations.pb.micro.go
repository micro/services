// Code generated by protoc-gen-micro. DO NOT EDIT.
// source: proto/locations.proto

package locations

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	_ "github.com/golang/protobuf/ptypes/timestamp"
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

// Api Endpoints for Locations service

func NewLocationsEndpoints() []*api.Endpoint {
	return []*api.Endpoint{}
}

// Client API for Locations service

type LocationsService interface {
	// Save a set of locations
	Save(ctx context.Context, in *SaveRequest, opts ...client.CallOption) (*SaveResponse, error)
	// Last locations for a set of users
	Last(ctx context.Context, in *LastRequest, opts ...client.CallOption) (*ListResponse, error)
	// Near returns the locations near a point at a given time
	Near(ctx context.Context, in *NearRequest, opts ...client.CallOption) (*ListResponse, error)
	// Read locations for a group of users between two points in time
	Read(ctx context.Context, in *ReadRequest, opts ...client.CallOption) (*ListResponse, error)
}

type locationsService struct {
	c    client.Client
	name string
}

func NewLocationsService(name string, c client.Client) LocationsService {
	return &locationsService{
		c:    c,
		name: name,
	}
}

func (c *locationsService) Save(ctx context.Context, in *SaveRequest, opts ...client.CallOption) (*SaveResponse, error) {
	req := c.c.NewRequest(c.name, "Locations.Save", in)
	out := new(SaveResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *locationsService) Last(ctx context.Context, in *LastRequest, opts ...client.CallOption) (*ListResponse, error) {
	req := c.c.NewRequest(c.name, "Locations.Last", in)
	out := new(ListResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *locationsService) Near(ctx context.Context, in *NearRequest, opts ...client.CallOption) (*ListResponse, error) {
	req := c.c.NewRequest(c.name, "Locations.Near", in)
	out := new(ListResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *locationsService) Read(ctx context.Context, in *ReadRequest, opts ...client.CallOption) (*ListResponse, error) {
	req := c.c.NewRequest(c.name, "Locations.Read", in)
	out := new(ListResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for Locations service

type LocationsHandler interface {
	// Save a set of locations
	Save(context.Context, *SaveRequest, *SaveResponse) error
	// Last locations for a set of users
	Last(context.Context, *LastRequest, *ListResponse) error
	// Near returns the locations near a point at a given time
	Near(context.Context, *NearRequest, *ListResponse) error
	// Read locations for a group of users between two points in time
	Read(context.Context, *ReadRequest, *ListResponse) error
}

func RegisterLocationsHandler(s server.Server, hdlr LocationsHandler, opts ...server.HandlerOption) error {
	type locations interface {
		Save(ctx context.Context, in *SaveRequest, out *SaveResponse) error
		Last(ctx context.Context, in *LastRequest, out *ListResponse) error
		Near(ctx context.Context, in *NearRequest, out *ListResponse) error
		Read(ctx context.Context, in *ReadRequest, out *ListResponse) error
	}
	type Locations struct {
		locations
	}
	h := &locationsHandler{hdlr}
	return s.Handle(s.NewHandler(&Locations{h}, opts...))
}

type locationsHandler struct {
	LocationsHandler
}

func (h *locationsHandler) Save(ctx context.Context, in *SaveRequest, out *SaveResponse) error {
	return h.LocationsHandler.Save(ctx, in, out)
}

func (h *locationsHandler) Last(ctx context.Context, in *LastRequest, out *ListResponse) error {
	return h.LocationsHandler.Last(ctx, in, out)
}

func (h *locationsHandler) Near(ctx context.Context, in *NearRequest, out *ListResponse) error {
	return h.LocationsHandler.Near(ctx, in, out)
}

func (h *locationsHandler) Read(ctx context.Context, in *ReadRequest, out *ListResponse) error {
	return h.LocationsHandler.Read(ctx, in, out)
}
