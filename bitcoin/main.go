package main

import (
	"github.com/micro/services/bitcoin/handler"
	pb "github.com/micro/services/bitcoin/proto"
	"micro.dev/v4/service"
	"micro.dev/v4/service/logger"
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
