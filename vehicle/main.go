package main

import (
	"github.com/micro/services/vehicle/handler"
	pb "github.com/micro/services/vehicle/proto"

	"github.com/micro/micro/v3/service"
	"github.com/micro/micro/v3/service/config"
	"github.com/micro/micro/v3/service/logger"
)

func main() {
	// Create service
	srv := service.New(
		service.Name("vehicle"),
		service.Version("latest"),
	)

	v, err := config.Get("dvla.api_key")
	if err != nil {
		logger.Fatalf("sunnha.api_key config not found: %v", err)
	}
	key := v.String("")
	if len(key) == 0 {
		logger.Fatal("dvla.api_key config not found")
	}

	h := handler.New(key)
	// Register handler
	pb.RegisterVehicleHandler(srv.Server(), h)
	pb.RegisterVehicleAdminHandler(srv.Server(), h)

	// Run service
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
