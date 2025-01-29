package main

import (
	"github.com/micro/services/holidays/handler"
	pb "github.com/micro/services/holidays/proto"

	"github.com/micro/micro/v5/service"
	"github.com/micro/micro/v5/service/logger"
)

func main() {
	// Create service
	srv := service.New(
		service.Name("holidays"),
		service.Version("latest"),
	)

	// Register handler
	pb.RegisterHolidaysHandler(srv.Server(), handler.New())

	// Run service
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
