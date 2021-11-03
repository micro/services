package main

import (
	"github.com/micro/micro/v3/service"
	"github.com/micro/micro/v3/service/logger"
	"github.com/micro/services/event/handler"
	pb "github.com/micro/services/event/proto"
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
