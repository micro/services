package main

import (
	admin "github.com/micro/services/pkg/service/proto"
	"github.com/micro/services/url/handler"
	pb "github.com/micro/services/url/proto"

	"micro.dev/v4/service"
	"micro.dev/v4/service/logger"
)

func main() {
	// Create service
	srv := service.New(
		service.Name("url"),
		service.Version("latest"),
	)
	h := handler.NewUrl(srv)
	// Register handler
	pb.RegisterUrlHandler(srv.Server(), h)
	admin.RegisterAdminHandler(srv.Server(), h)

	// Run service
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
