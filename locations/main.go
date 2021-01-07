package main

import (
	"github.com/micro/services/locations/handler"
	pb "github.com/micro/services/locations/proto"

	"github.com/micro/micro/v3/service"
	"github.com/micro/micro/v3/service/logger"
)

func main() {
	// Create service
	srv := service.New(
		service.Name("locations"),
		service.Version("latest"),
	)

	// Register handler
	pb.RegisterLocationsHandler(srv.Server(), new(handler.Locations))

	// Run service
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
