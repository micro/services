// Code generated by protoc-gen-micro. DO NOT EDIT.
// source: proto/users.proto

package users

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

// Api Endpoints for Users service

func NewUsersEndpoints() []*api.Endpoint {
	return []*api.Endpoint{}
}

// Client API for Users service

type UsersService interface {
	Create(ctx context.Context, in *CreateRequest, opts ...client.CallOption) (*CreateResponse, error)
	Read(ctx context.Context, in *ReadRequest, opts ...client.CallOption) (*ReadResponse, error)
	ReadByEmail(ctx context.Context, in *ReadByEmailRequest, opts ...client.CallOption) (*ReadByEmailResponse, error)
	Update(ctx context.Context, in *UpdateRequest, opts ...client.CallOption) (*UpdateResponse, error)
	Delete(ctx context.Context, in *DeleteRequest, opts ...client.CallOption) (*DeleteResponse, error)
	List(ctx context.Context, in *ListRequest, opts ...client.CallOption) (*ListResponse, error)
	// Login using email and password returns the users profile and a token
	Login(ctx context.Context, in *LoginRequest, opts ...client.CallOption) (*LoginResponse, error)
	// Logout expires all tokens for the user
	Logout(ctx context.Context, in *LogoutRequest, opts ...client.CallOption) (*LogoutResponse, error)
	// Validate a token, each time a token is validated it extends its lifetime for another week
	Validate(ctx context.Context, in *ValidateRequest, opts ...client.CallOption) (*ValidateResponse, error)
}

type usersService struct {
	c    client.Client
	name string
}

func NewUsersService(name string, c client.Client) UsersService {
	return &usersService{
		c:    c,
		name: name,
	}
}

func (c *usersService) Create(ctx context.Context, in *CreateRequest, opts ...client.CallOption) (*CreateResponse, error) {
	req := c.c.NewRequest(c.name, "Users.Create", in)
	out := new(CreateResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *usersService) Read(ctx context.Context, in *ReadRequest, opts ...client.CallOption) (*ReadResponse, error) {
	req := c.c.NewRequest(c.name, "Users.Read", in)
	out := new(ReadResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *usersService) ReadByEmail(ctx context.Context, in *ReadByEmailRequest, opts ...client.CallOption) (*ReadByEmailResponse, error) {
	req := c.c.NewRequest(c.name, "Users.ReadByEmail", in)
	out := new(ReadByEmailResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *usersService) Update(ctx context.Context, in *UpdateRequest, opts ...client.CallOption) (*UpdateResponse, error) {
	req := c.c.NewRequest(c.name, "Users.Update", in)
	out := new(UpdateResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *usersService) Delete(ctx context.Context, in *DeleteRequest, opts ...client.CallOption) (*DeleteResponse, error) {
	req := c.c.NewRequest(c.name, "Users.Delete", in)
	out := new(DeleteResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *usersService) List(ctx context.Context, in *ListRequest, opts ...client.CallOption) (*ListResponse, error) {
	req := c.c.NewRequest(c.name, "Users.List", in)
	out := new(ListResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *usersService) Login(ctx context.Context, in *LoginRequest, opts ...client.CallOption) (*LoginResponse, error) {
	req := c.c.NewRequest(c.name, "Users.Login", in)
	out := new(LoginResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *usersService) Logout(ctx context.Context, in *LogoutRequest, opts ...client.CallOption) (*LogoutResponse, error) {
	req := c.c.NewRequest(c.name, "Users.Logout", in)
	out := new(LogoutResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *usersService) Validate(ctx context.Context, in *ValidateRequest, opts ...client.CallOption) (*ValidateResponse, error) {
	req := c.c.NewRequest(c.name, "Users.Validate", in)
	out := new(ValidateResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for Users service

type UsersHandler interface {
	Create(context.Context, *CreateRequest, *CreateResponse) error
	Read(context.Context, *ReadRequest, *ReadResponse) error
	ReadByEmail(context.Context, *ReadByEmailRequest, *ReadByEmailResponse) error
	Update(context.Context, *UpdateRequest, *UpdateResponse) error
	Delete(context.Context, *DeleteRequest, *DeleteResponse) error
	List(context.Context, *ListRequest, *ListResponse) error
	// Login using email and password returns the users profile and a token
	Login(context.Context, *LoginRequest, *LoginResponse) error
	// Logout expires all tokens for the user
	Logout(context.Context, *LogoutRequest, *LogoutResponse) error
	// Validate a token, each time a token is validated it extends its lifetime for another week
	Validate(context.Context, *ValidateRequest, *ValidateResponse) error
}

func RegisterUsersHandler(s server.Server, hdlr UsersHandler, opts ...server.HandlerOption) error {
	type users interface {
		Create(ctx context.Context, in *CreateRequest, out *CreateResponse) error
		Read(ctx context.Context, in *ReadRequest, out *ReadResponse) error
		ReadByEmail(ctx context.Context, in *ReadByEmailRequest, out *ReadByEmailResponse) error
		Update(ctx context.Context, in *UpdateRequest, out *UpdateResponse) error
		Delete(ctx context.Context, in *DeleteRequest, out *DeleteResponse) error
		List(ctx context.Context, in *ListRequest, out *ListResponse) error
		Login(ctx context.Context, in *LoginRequest, out *LoginResponse) error
		Logout(ctx context.Context, in *LogoutRequest, out *LogoutResponse) error
		Validate(ctx context.Context, in *ValidateRequest, out *ValidateResponse) error
	}
	type Users struct {
		users
	}
	h := &usersHandler{hdlr}
	return s.Handle(s.NewHandler(&Users{h}, opts...))
}

type usersHandler struct {
	UsersHandler
}

func (h *usersHandler) Create(ctx context.Context, in *CreateRequest, out *CreateResponse) error {
	return h.UsersHandler.Create(ctx, in, out)
}

func (h *usersHandler) Read(ctx context.Context, in *ReadRequest, out *ReadResponse) error {
	return h.UsersHandler.Read(ctx, in, out)
}

func (h *usersHandler) ReadByEmail(ctx context.Context, in *ReadByEmailRequest, out *ReadByEmailResponse) error {
	return h.UsersHandler.ReadByEmail(ctx, in, out)
}

func (h *usersHandler) Update(ctx context.Context, in *UpdateRequest, out *UpdateResponse) error {
	return h.UsersHandler.Update(ctx, in, out)
}

func (h *usersHandler) Delete(ctx context.Context, in *DeleteRequest, out *DeleteResponse) error {
	return h.UsersHandler.Delete(ctx, in, out)
}

func (h *usersHandler) List(ctx context.Context, in *ListRequest, out *ListResponse) error {
	return h.UsersHandler.List(ctx, in, out)
}

func (h *usersHandler) Login(ctx context.Context, in *LoginRequest, out *LoginResponse) error {
	return h.UsersHandler.Login(ctx, in, out)
}

func (h *usersHandler) Logout(ctx context.Context, in *LogoutRequest, out *LogoutResponse) error {
	return h.UsersHandler.Logout(ctx, in, out)
}

func (h *usersHandler) Validate(ctx context.Context, in *ValidateRequest, out *ValidateResponse) error {
	return h.UsersHandler.Validate(ctx, in, out)
}
