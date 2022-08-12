package main

import (
	"github.com/micro/micro/v3/service"
	"github.com/micro/micro/v3/service/logger"
	"github.com/micro/services/ethereum/handler"
	pb "github.com/micro/services/ethereum/proto"
)

func main() {
	// Create service
	srv := service.New(
		service.Name("ethereum"),
	)

	// Register handler
	pb.RegisterEthereumHandler(srv.Server(), handler.New())

	// Run service
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
