package main

import (
	"github.com/micro/micro/v3/service"
	"github.com/micro/micro/v3/service/config"
	"github.com/micro/micro/v3/service/logger"
	"github.com/micro/services/sunnah/handler"
	pb "github.com/micro/services/sunnah/proto"
)

func main() {
	// Create service
	srv := service.New(
		service.Name("sunnah"),
		service.Version("latest"),
	)

	v, err := config.Get("sunnah.api_key")
	if err != nil {
		logger.Fatalf("sunnha.api_key config not found: %v", err)
	}
	key := v.String("")
	if len(key) == 0 {
		logger.Fatal("sunnah.api_key config not found")
	}

	// Register handler
	pb.RegisterSunnahHandler(srv.Server(), handler.New(key))

	// Run service
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
