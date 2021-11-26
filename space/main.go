package main

import (
	"github.com/micro/services/space/handler"
	pb "github.com/micro/services/space/proto"
	"github.com/micro/micro/v3/service"
	"github.com/micro/micro/v3/service/logger"
)

func main() {
	// Create service
	srv := service.New(
		service.Name("space"),
		service.Version("latest"),
	)

	// Register handler
	pb.RegisterSpaceHandler(srv.Server(), new(handler.Space))

	// Run service
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
