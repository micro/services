package main

import (
	"github.com/micro/services/event/handler"
	pb "github.com/micro/services/event/proto"
	"micro.dev/v4/service"
	"micro.dev/v4/service/logger"
)

func main() {
	// Create service
	srv := service.New(
		service.Name("event"),
		service.Version("latest"),
	)

	// Register handler
	pb.RegisterEventHandler(srv.Server(), new(handler.Event))

	// Run service
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
