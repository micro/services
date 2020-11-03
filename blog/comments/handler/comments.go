package handler

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/micro/dev/model"
	"github.com/micro/micro/v3/service/store"
	pb "github.com/micro/services/blog/comments/proto"
)

type Comments struct {
	comments  model.Model
	idIndex   model.Index
	postIndex model.Index
}

func NewComments() *Comments {
	postIndex := model.ByEquality("post")
	postIndex.Order.Type = model.OrderTypeDesc
	postIndex.Order.FieldName = "created"

	idIndex := model.ByEquality("id")
	idIndex.Order.Type = model.OrderTypeUnordered

	return &Comments{
		comments:  model.New(store.DefaultStore, "users", model.Indexes(postIndex), nil),
		postIndex: postIndex,
		idIndex:   idIndex,
	}
}

func (c *Comments) New(ctx context.Context, req *pb.NewRequest, rsp *pb.NewResponse) error {
	return c.comments.Save(pb.Comment{
		Id:      uuid.New().String(),
		Post:    req.Post,
		Author:  req.Author,
		Message: req.Message,
		Created: time.Now().Unix(),
	})
}

func (c *Comments) List(ctx context.Context, req *pb.ListRequest, rsp *pb.ListResponse) error {
	return c.comments.List(c.postIndex.ToQuery(req.Post), &rsp.Comments)
}
