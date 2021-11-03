package main

import (
	"github.com/micro/micro/v3/service"
	"github.com/micro/micro/v3/service/logger"
	"github.com/micro/services/stream/handler"
	pb "github.com/micro/services/stream/proto"
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
