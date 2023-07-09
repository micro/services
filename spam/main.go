package main

import (
	"github.com/micro/services/spam/handler"
	pb "github.com/micro/services/spam/proto"

	"micro.dev/v4/service"
	"micro.dev/v4/service/logger"
)

func main() {
	// Create service
	srv := service.New(
		service.Name("spam"),
		service.Version("latest"),
	)

	// Register handler
	pb.RegisterSpamHandler(srv.Server(), handler.New())

	// Run service
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
