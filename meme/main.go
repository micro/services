package main

import (
	"github.com/micro/micro/v3/service"
	"github.com/micro/micro/v3/service/logger"
	"github.com/micro/services/meme/handler"
	pb "github.com/micro/services/meme/proto"
)

func main() {
	// Create service
	srv := service.New(
		service.Name("meme"),
		service.Version("latest"),
	)

	// Register handler
	pb.RegisterMemeHandler(srv.Server(), handler.New())

	// Run service
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
