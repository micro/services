package main

import (
	"github.com/micro/services/weather/handler"
	pb "github.com/micro/services/weather/proto"
	"github.com/micro/micro/v3/service"
	"github.com/micro/micro/v3/service/logger"
)

func main() {
	// Create service
	srv := service.New(
		service.Name("weather"),
		service.Version("latest"),
	)

	// Register handler
	pb.RegisterWeatherHandler(srv.Server(), handler.New())

	// Run service
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
