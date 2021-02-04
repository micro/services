// Code generated by protoc-gen-micro. DO NOT EDIT.
// source: proto/invites.proto

package invites

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	_ "github.com/golang/protobuf/ptypes/wrappers"
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

// Api Endpoints for Invites service

func NewInvitesEndpoints() []*api.Endpoint {
	return []*api.Endpoint{}
}

// Client API for Invites service

type InvitesService interface {
	// Create an invite
	Create(ctx context.Context, in *CreateRequest, opts ...client.CallOption) (*CreateResponse, error)
	// Read an invite using ID or code
	Read(ctx context.Context, in *ReadRequest, opts ...client.CallOption) (*ReadResponse, error)
	// List invited for a group or specific email
	List(ctx context.Context, in *ListRequest, opts ...client.CallOption) (*ListResponse, error)
	// Delete an invite
	Delete(ctx context.Context, in *DeleteRequest, opts ...client.CallOption) (*DeleteResponse, error)
}

type invitesService struct {
	c    client.Client
	name string
}

func NewInvitesService(name string, c client.Client) InvitesService {
	return &invitesService{
		c:    c,
		name: name,
	}
}

func (c *invitesService) Create(ctx context.Context, in *CreateRequest, opts ...client.CallOption) (*CreateResponse, error) {
	req := c.c.NewRequest(c.name, "Invites.Create", in)
	out := new(CreateResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *invitesService) Read(ctx context.Context, in *ReadRequest, opts ...client.CallOption) (*ReadResponse, error) {
	req := c.c.NewRequest(c.name, "Invites.Read", in)
	out := new(ReadResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *invitesService) List(ctx context.Context, in *ListRequest, opts ...client.CallOption) (*ListResponse, error) {
	req := c.c.NewRequest(c.name, "Invites.List", in)
	out := new(ListResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *invitesService) Delete(ctx context.Context, in *DeleteRequest, opts ...client.CallOption) (*DeleteResponse, error) {
	req := c.c.NewRequest(c.name, "Invites.Delete", in)
	out := new(DeleteResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for Invites service

type InvitesHandler interface {
	// Create an invite
	Create(context.Context, *CreateRequest, *CreateResponse) error
	// Read an invite using ID or code
	Read(context.Context, *ReadRequest, *ReadResponse) error
	// List invited for a group or specific email
	List(context.Context, *ListRequest, *ListResponse) error
	// Delete an invite
	Delete(context.Context, *DeleteRequest, *DeleteResponse) error
}

func RegisterInvitesHandler(s server.Server, hdlr InvitesHandler, opts ...server.HandlerOption) error {
	type invites interface {
		Create(ctx context.Context, in *CreateRequest, out *CreateResponse) error
		Read(ctx context.Context, in *ReadRequest, out *ReadResponse) error
		List(ctx context.Context, in *ListRequest, out *ListResponse) error
		Delete(ctx context.Context, in *DeleteRequest, out *DeleteResponse) error
	}
	type Invites struct {
		invites
	}
	h := &invitesHandler{hdlr}
	return s.Handle(s.NewHandler(&Invites{h}, opts...))
}

type invitesHandler struct {
	InvitesHandler
}

func (h *invitesHandler) Create(ctx context.Context, in *CreateRequest, out *CreateResponse) error {
	return h.InvitesHandler.Create(ctx, in, out)
}

func (h *invitesHandler) Read(ctx context.Context, in *ReadRequest, out *ReadResponse) error {
	return h.InvitesHandler.Read(ctx, in, out)
}

func (h *invitesHandler) List(ctx context.Context, in *ListRequest, out *ListResponse) error {
	return h.InvitesHandler.List(ctx, in, out)
}

func (h *invitesHandler) Delete(ctx context.Context, in *DeleteRequest, out *DeleteResponse) error {
	return h.InvitesHandler.Delete(ctx, in, out)
}
