package main

import (
	"github.com/micro/services/gifs/handler"
	pb "github.com/micro/services/gifs/proto"

	"micro.dev/v4/service"
	"micro.dev/v4/service/logger"
)

func main() {
	// Create service
	srv := service.New(
		service.Name("gifs"),
		service.Version("latest"),
	)

	// Register handler
	pb.RegisterGifsHandler(srv.Server(), handler.New())

	// Run service
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
