package handler

import (
	"context"

	proto "github.com/micro/services/blog/proto"
	comments "github.com/micro/services/comments/proto"
	posts "github.com/micro/services/posts/proto"
	tags "github.com/micro/services/tags/proto"
)

type Blog struct {
	ps posts.PostsService
	cs comments.CommentsService
	ts tags.TagsService
}

func NewBlog(ps posts.PostsService,
	cs comments.CommentsService,
	ts tags.TagsService) *Blog {
	return &Blog{
		ps: ps,
		cs: cs,
		ts: ts,
	}
}

func (e *Blog) Latest(ctx context.Context, req *proto.LatestRequest, rsp *proto.LatestResponse) error {
	return nil
}
