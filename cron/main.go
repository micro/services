package main

import (
	"github.com/micro/services/cron/handler"
	pb "github.com/micro/services/cron/proto"

	"github.com/micro/micro/v3/service"
	"github.com/micro/micro/v3/service/logger"
)

func main() {
	// Create service
	srv := service.New(
		service.Name("cron"),
	)

	// Register handler
	pb.RegisterCronHandler(srv.Server(), handler.New())

	// Run service
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
