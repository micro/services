package main

import (
	"github.com/micro/micro/v5/service"
	"github.com/micro/micro/v5/service/logger"
	"github.com/micro/services/carbon/handler"
	pb "github.com/micro/services/carbon/proto"
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
