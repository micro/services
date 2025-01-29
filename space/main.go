package main

import (
	"github.com/micro/micro/v5/service"
	"github.com/micro/micro/v5/service/logger"
	admin "github.com/micro/services/pkg/service/proto"
	"github.com/micro/services/space/handler"
	pb "github.com/micro/services/space/proto"
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
