package main

import (
	"github.com/micro/services/app/handler"
	pb "github.com/micro/services/app/proto"
	admin "github.com/micro/services/pkg/service/proto"
	"micro.dev/v4/service"
	"micro.dev/v4/service/logger"
)

func main() {
	// Create service
	srv := service.New(
		service.Name("app"),
		service.Version("latest"),
	)

	h := handler.New(srv)
	// Register handler
	pb.RegisterAppHandler(srv.Server(), h)
	admin.RegisterAdminHandler(srv.Server(), h)

	// Run service
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
