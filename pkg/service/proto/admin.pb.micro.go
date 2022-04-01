// Code generated by protoc-gen-micro. DO NOT EDIT.
// source: proto/admin.proto

package admin

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

// Api Endpoints for Admin service

func NewAdminEndpoints() []*api.Endpoint {
	return []*api.Endpoint{}
}

// Client API for Admin service

type AdminService interface {
	DeleteData(ctx context.Context, in *DeleteDataRequest, opts ...client.CallOption) (*DeleteDataResponse, error)
	Usage(ctx context.Context, in *UsageRequest, opts ...client.CallOption) (*UsageResponse, error)
}

type adminService struct {
	c    client.Client
	name string
}

func NewAdminService(name string, c client.Client) AdminService {
	return &adminService{
		c:    c,
		name: name,
	}
}

func (c *adminService) DeleteData(ctx context.Context, in *DeleteDataRequest, opts ...client.CallOption) (*DeleteDataResponse, error) {
	req := c.c.NewRequest(c.name, "Admin.DeleteData", in)
	out := new(DeleteDataResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *adminService) Usage(ctx context.Context, in *UsageRequest, opts ...client.CallOption) (*UsageResponse, error) {
	req := c.c.NewRequest(c.name, "Admin.Usage", in)
	out := new(UsageResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for Admin service

type AdminHandler interface {
	DeleteData(context.Context, *DeleteDataRequest, *DeleteDataResponse) error
	Usage(context.Context, *UsageRequest, *UsageResponse) error
}

func RegisterAdminHandler(s server.Server, hdlr AdminHandler, opts ...server.HandlerOption) error {
	type admin interface {
		DeleteData(ctx context.Context, in *DeleteDataRequest, out *DeleteDataResponse) error
		Usage(ctx context.Context, in *UsageRequest, out *UsageResponse) error
	}
	type Admin struct {
		admin
	}
	h := &adminHandler{hdlr}
	return s.Handle(s.NewHandler(&Admin{h}, opts...))
}

type adminHandler struct {
	AdminHandler
}

func (h *adminHandler) DeleteData(ctx context.Context, in *DeleteDataRequest, out *DeleteDataResponse) error {
	return h.AdminHandler.DeleteData(ctx, in, out)
}

func (h *adminHandler) Usage(ctx context.Context, in *UsageRequest, out *UsageResponse) error {
	return h.AdminHandler.Usage(ctx, in, out)
}
