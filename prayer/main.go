package main

import (
	"github.com/micro/services/prayer/handler"
	pb "github.com/micro/services/prayer/proto"
	"micro.dev/v4/service"
	"micro.dev/v4/service/logger"
)

func main() {
	// Create service
	srv := service.New(
		service.Name("prayer"),
		service.Version("latest"),
	)

	// Register handler
	pb.RegisterPrayerHandler(srv.Server(), handler.New(srv.Client()))

	// Run service
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
