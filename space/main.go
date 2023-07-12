package main

import (
	admin "github.com/micro/services/pkg/service/proto"
	"github.com/micro/services/space/handler"
	pb "github.com/micro/services/space/proto"
	"micro.dev/v4/service"
	"micro.dev/v4/service/logger"
)

func main() {
	// Create service
	srv := service.New(
		service.Name("space"),
		service.Version("latest"),
	)

	h := handler.NewSpace(srv)
	// Register handler
	admin.RegisterAdminHandler(srv.Server(), h)
	pb.RegisterSpaceHandler(srv.Server(), h)

	// Run service
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
