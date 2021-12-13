package main

import (
	"github.com/micro/micro/v3/service"
	"github.com/micro/micro/v3/service/logger"
	"github.com/micro/services/app/handler"
	pb "github.com/micro/services/app/proto"
)

func main() {
	// Create service
	srv := service.New(
		service.Name("app"),
		service.Version("latest"),
	)

	// Register handler
	pb.RegisterAppHandler(srv.Server(), handler.New())

	// Run service
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
