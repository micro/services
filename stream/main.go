package main

import (
	"github.com/micro/services/stream/handler"
	pb "github.com/micro/services/stream/proto"
	"micro.dev/v4/service"
	"micro.dev/v4/service/logger"
)

func main() {
	// Create service
	srv := service.New(
		service.Name("stream"),
		service.Version("latest"),
	)

	// Register handler
	pb.RegisterStreamHandler(srv.Server(), handler.New())

	// Run service
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
