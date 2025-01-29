package main

import (
	"github.com/micro/services/cache/handler"
	pb "github.com/micro/services/cache/proto"
	adminpb "github.com/micro/services/pkg/service/proto"

	"github.com/micro/micro/v5/service"
	"github.com/micro/micro/v5/service/logger"
)

func main() {
	// Create service
	srv := service.New(
		service.Name("cache"),
		service.Version("latest"),
	)

	// Register handler
	c := new(handler.Cache)
	pb.RegisterCacheHandler(srv.Server(), c)
	adminpb.RegisterAdminHandler(srv.Server(), c)

	// Run service
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
