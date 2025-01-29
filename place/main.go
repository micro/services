package main

import (
	"github.com/micro/services/place/handler"
	pb "github.com/micro/services/place/proto"

	"github.com/micro/micro/v5/service"
	"github.com/micro/micro/v5/service/logger"
)

func main() {
	// Create service
	srv := service.New(
		service.Name("place"),
		service.Version("latest"),
	)

	// Register handler
	pb.RegisterPlaceHandler(srv.Server(), handler.New())

	// Run service
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
