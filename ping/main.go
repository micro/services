package main

import (
	"github.com/micro/services/ping/handler"
	pb "github.com/micro/services/ping/proto"
	"micro.dev/v4/service"
	"micro.dev/v4/service/logger"
)

func main() {
	// Create service
	srv := service.New(
		service.Name("ping"),
		service.Version("latest"),
	)

	// Register handler
	pb.RegisterPingHandler(srv.Server(), new(handler.Ping))

	// Run service
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
