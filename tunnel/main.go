package main

import (
	"github.com/micro/services/tunnel/handler"
	pb "github.com/micro/services/tunnel/proto"
	"micro.dev/v4/service"
	"micro.dev/v4/service/logger"
)

func main() {
	// Create service
	srv := service.New(
		service.Name("tunnel"),
		service.Version("latest"),
	)

	// Register handler
	pb.RegisterTunnelHandler(srv.Server(), handler.New())

	// Run service
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
