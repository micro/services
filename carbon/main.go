package main

import (
	"github.com/micro/services/carbon/handler"
	pb "github.com/micro/services/carbon/proto"
	"micro.dev/v4/service"
	"micro.dev/v4/service/logger"
)

func main() {
	// Create service
	srv := service.New(
		service.Name("carbon"),
		service.Version("latest"),
	)

	// Register handler
	pb.RegisterCarbonHandler(srv.Server(), handler.New())

	// Run service
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
