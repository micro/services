package main

import (
	admin "github.com/micro/services/pkg/service/proto"
	"github.com/micro/services/search/handler"
	pb "github.com/micro/services/search/proto"
	"micro.dev/v4/service"
	"micro.dev/v4/service/logger"
)

func main() {
	// Create service
	srv := service.New(
		service.Name("search"),
		service.Version("latest"),
	)

	h := handler.New(srv)
	// Register handler
	pb.RegisterSearchHandler(srv.Server(), h)
	admin.RegisterAdminHandler(srv.Server(), h)

	// Run service
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
