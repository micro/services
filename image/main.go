package main

import (
	"github.com/micro/services/image/handler"
	pb "github.com/micro/services/image/proto"

	"github.com/micro/micro/v3/service"
	"github.com/micro/micro/v3/service/logger"
)

func main() {
	// Create service
	srv := service.New(
		service.Name("image"),
		service.Version("latest"),
	)

	// Register handler
	pb.RegisterImageHandler(srv.Server(), handler.NewImage())

	// Run service
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
