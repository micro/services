package main

import (
	"github.com/micro/services/streams/handler"
	pb "github.com/micro/services/streams/proto"

	"github.com/micro/micro/v3/service"
	"github.com/micro/micro/v3/service/logger"
)

func main() {
	// Create service
	srv := service.New(
		service.Name("streams"),
		service.Version("latest"),
	)

	// Register handler
	pb.RegisterStreamsHandler(srv.Server(), new(handler.Streams))

	// Run service
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
