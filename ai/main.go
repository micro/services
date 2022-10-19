package main

import (
	"github.com/micro/services/ai/handler"
	pb "github.com/micro/services/ai/proto"

	"github.com/micro/micro/v3/service"
	"github.com/micro/micro/v3/service/logger"
)

func main() {
	// Create service
	srv := service.New(
		service.Name("ai"),
	)

	// Register handler
	pb.RegisterAiHandler(srv.Server(), handler.New())

	// Run service
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
