package main

import (
	"github.com/micro/services/password/handler"
	pb "github.com/micro/services/password/proto"

	"github.com/micro/micro/v5/service"
	"github.com/micro/micro/v5/service/logger"
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
