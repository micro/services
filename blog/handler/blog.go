package handler

import (
	"context"

	log "github.com/micro/micro/v3/service/logger"

	blog "github.com/micro/services/blog/proto"
)

type Blog struct{}

func (e *Blog) Latest(ctx context.Context, req *blog.LatestRequest, rsp *blog.LatestResponse) error {
	log.Info("Received Blog.Call request")
	rsp.Msg = "Hello " + req.Name
	return nil
}
