package main

import (
	"github.com/micro/services/image/handler"
	pb "github.com/micro/services/image/proto"
	admin "github.com/micro/services/pkg/service/proto"

	"micro.dev/v4/service"
	"micro.dev/v4/service/logger"
)

func main() {
	// Create service
	srv := service.New(
		service.Name("image"),
		service.Version("latest"),
	)

	h := handler.NewImage()
	// Register handler
	pb.RegisterImageHandler(srv.Server(), h)
	admin.RegisterAdminHandler(srv.Server(), h)

	// Run service
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
