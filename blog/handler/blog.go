package handler

import (
	"context"

	blog "github.com/micro/services/blog/proto"
)

type Blog struct{}

func (e *Blog) Latest(ctx context.Context, req *blog.LatestRequest, rsp *blog.LatestResponse) error {
	return nil
}
