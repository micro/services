package main

import (
	"github.com/micro/micro/v3/service"
	"github.com/micro/micro/v3/service/logger"
	"github.com/micro/services/forex/handler"
	pb "github.com/micro/services/forex/proto"
)

func main() {
	// Create service
	srv := service.New(
		service.Name("forex"),
		service.Version("latest"),
	)

	// Register handler
	pb.RegisterForexHandler(srv.Server(), handler.New())

	// Run service
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
