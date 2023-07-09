package main

import (
	"github.com/micro/services/time/handler"
	pb "github.com/micro/services/time/proto"
	"micro.dev/v4/service"
	"micro.dev/v4/service/logger"
)

func main() {
	// Create service
	srv := service.New(
		service.Name("time"),
		service.Version("latest"),
	)

	// Register handler
	pb.RegisterTimeHandler(srv.Server(), handler.New())

	// Run service
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
