package main

import (
	"github.com/micro/services/files/handler"
	pb "github.com/micro/services/files/proto"

	"github.com/micro/micro/v3/service"
	"github.com/micro/micro/v3/service/logger"
)

func main() {
	// Create service
	srv := service.New(
		service.Name("files"),
		service.Version("latest"),
	)

	// Register handler
	pb.RegisterFilesHandler(srv.Server(), new(handler.Files))

	// Run service
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
