package handler

import (
	"context"

	"github.com/micro/micro/v3/service/logger"
	pb "github.com/micro/services/blog/comments/proto"
)

type Comments struct{}

// Call is a single request handler called via client.Call or the generated client code
func (c *Comments) Save(ctx context.Context, req *pb.Request, rsp *pb.Response) error {
	logger.Info("Not yet implemented")
	return nil
}
