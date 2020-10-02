package handler

import (
	"context"

	"github.com/micro/micro/v3/service/logger"
	pb "github.com/micro/services/blog/search/proto"
)

type Search struct{}

func (s *Search) Index(ctx context.Context, req *pb.IndexRequest, rsp *pb.IndexResponse) error {
	logger.Info("Received Search.Index request")
	return nil
}

func (s *Search) Search(ctx context.Context, req *pb.SearchRequest, rsp *pb.SearchResponse) error {
	logger.Info("Received Search.Search request")
	return nil
}
