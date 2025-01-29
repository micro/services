package main

import (
	"github.com/micro/micro/v5/service"
	"github.com/micro/micro/v5/service/logger"
	"github.com/micro/services/bitcoin/handler"
	pb "github.com/micro/services/bitcoin/proto"
)

func main() {
	// Create service
	srv := service.New(
		service.Name("bitcoin"),
		service.Version("latest"),
	)

	// Register handler
	pb.RegisterBitcoinHandler(srv.Server(), handler.New())

	// Run service
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
