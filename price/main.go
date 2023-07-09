package main

import (
	"github.com/micro/services/price/handler"
	pb "github.com/micro/services/price/proto"
	"micro.dev/v4/service"
	"micro.dev/v4/service/logger"
)

func main() {
	// Create service
	srv := service.New(
		service.Name("price"),
		service.Version("latest"),
	)

	// Register handler
	pb.RegisterPriceHandler(srv.Server(), handler.New())

	// Run service
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
