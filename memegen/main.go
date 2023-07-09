package main

import (
	"github.com/micro/services/memegen/handler"
	pb "github.com/micro/services/memegen/proto"
	"micro.dev/v4/service"
	"micro.dev/v4/service/logger"
)

func main() {
	// Create service
	srv := service.New(
		service.Name("memegen"),
		service.Version("latest"),
	)

	// Register handler
	pb.RegisterMemegenHandler(srv.Server(), handler.New())

	// Run service
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
