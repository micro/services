package main

import (
	"github.com/micro/services/secret/handler"
	pb "github.com/micro/services/secret/proto"

	"github.com/micro/micro/v5/service"
	"github.com/micro/micro/v5/service/logger"
)

func main() {
	// Create service
	srv := service.New(
		service.Name("secret"),
	)

	// Register handler
	pb.RegisterSecretHandler(srv.Server(), handler.New())

	// Run service
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
