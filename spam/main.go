package main

import (
	"github.com/micro/services/spam/handler"
	pb "github.com/micro/services/spam/proto"

	"github.com/micro/micro/v3/service"
	"github.com/micro/micro/v3/service/logger"
)

func main() {
	// Create service
	srv := service.New(
		service.Name("spam"),
		service.Version("latest"),
	)

	// Register handler
	pb.RegisterSpamHandler(srv.Server(), new(handler.Spam))

	// Run service
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
