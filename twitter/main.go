package main

import (
	"github.com/micro/services/twitter/handler"
	pb "github.com/micro/services/twitter/proto"
	"micro.dev/v4/service"
	"micro.dev/v4/service/logger"
)

func main() {
	// Create service
	srv := service.New(
		service.Name("twitter"),
		service.Version("latest"),
	)

	// Register handler
	pb.RegisterTwitterHandler(srv.Server(), handler.New())

	// Run service
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
