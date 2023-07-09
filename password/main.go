package main

import (
	"github.com/micro/services/password/handler"
	pb "github.com/micro/services/password/proto"

	"micro.dev/v4/service"
	"micro.dev/v4/service/logger"
)

func main() {
	// Create service
	srv := service.New(
		service.Name("password"),
		service.Version("latest"),
	)

	// Register handler
	pb.RegisterPasswordHandler(srv.Server(), new(handler.Password))

	// Run service
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
