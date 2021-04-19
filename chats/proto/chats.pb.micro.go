// Code generated by protoc-gen-micro. DO NOT EDIT.
// source: proto/chats.proto

package chats

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

// Api Endpoints for Chats service

func NewChatsEndpoints() []*api.Endpoint {
	return []*api.Endpoint{}
}

// Client API for Chats service

type ChatsService interface {
	// Create a chat between two or more users, if a chat already exists for these users, the existing
	// chat will be returned
	CreateChat(ctx context.Context, in *CreateChatRequest, opts ...client.CallOption) (*CreateChatResponse, error)
	// Create a message within a chat
	CreateMessage(ctx context.Context, in *CreateMessageRequest, opts ...client.CallOption) (*CreateMessageResponse, error)
	// List the messages within a chat in reverse chronological order, using sent_before to
	// offset as older messages need to be loaded
	ListMessages(ctx context.Context, in *ListMessagesRequest, opts ...client.CallOption) (*ListMessagesResponse, error)
}

type chatsService struct {
	c    client.Client
	name string
}

func NewChatsService(name string, c client.Client) ChatsService {
	return &chatsService{
		c:    c,
		name: name,
	}
}

func (c *chatsService) CreateChat(ctx context.Context, in *CreateChatRequest, opts ...client.CallOption) (*CreateChatResponse, error) {
	req := c.c.NewRequest(c.name, "Chats.CreateChat", in)
	out := new(CreateChatResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *chatsService) CreateMessage(ctx context.Context, in *CreateMessageRequest, opts ...client.CallOption) (*CreateMessageResponse, error) {
	req := c.c.NewRequest(c.name, "Chats.CreateMessage", in)
	out := new(CreateMessageResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *chatsService) ListMessages(ctx context.Context, in *ListMessagesRequest, opts ...client.CallOption) (*ListMessagesResponse, error) {
	req := c.c.NewRequest(c.name, "Chats.ListMessages", in)
	out := new(ListMessagesResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for Chats service

type ChatsHandler interface {
	// Create a chat between two or more users, if a chat already exists for these users, the existing
	// chat will be returned
	CreateChat(context.Context, *CreateChatRequest, *CreateChatResponse) error
	// Create a message within a chat
	CreateMessage(context.Context, *CreateMessageRequest, *CreateMessageResponse) error
	// List the messages within a chat in reverse chronological order, using sent_before to
	// offset as older messages need to be loaded
	ListMessages(context.Context, *ListMessagesRequest, *ListMessagesResponse) error
}

func RegisterChatsHandler(s server.Server, hdlr ChatsHandler, opts ...server.HandlerOption) error {
	type chats interface {
		CreateChat(ctx context.Context, in *CreateChatRequest, out *CreateChatResponse) error
		CreateMessage(ctx context.Context, in *CreateMessageRequest, out *CreateMessageResponse) error
		ListMessages(ctx context.Context, in *ListMessagesRequest, out *ListMessagesResponse) error
	}
	type Chats struct {
		chats
	}
	h := &chatsHandler{hdlr}
	return s.Handle(s.NewHandler(&Chats{h}, opts...))
}

type chatsHandler struct {
	ChatsHandler
}

func (h *chatsHandler) CreateChat(ctx context.Context, in *CreateChatRequest, out *CreateChatResponse) error {
	return h.ChatsHandler.CreateChat(ctx, in, out)
}

func (h *chatsHandler) CreateMessage(ctx context.Context, in *CreateMessageRequest, out *CreateMessageResponse) error {
	return h.ChatsHandler.CreateMessage(ctx, in, out)
}

func (h *chatsHandler) ListMessages(ctx context.Context, in *ListMessagesRequest, out *ListMessagesResponse) error {
	return h.ChatsHandler.ListMessages(ctx, in, out)
}
