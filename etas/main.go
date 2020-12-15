package main

import (
	"etas/handler"
	pb "etas/proto"

	"googlemaps.github.io/maps"

	"github.com/micro/micro/v3/service"
	"github.com/micro/micro/v3/service/config"
	"github.com/micro/micro/v3/service/logger"
)

func main() {
	// Create service
	srv := service.New(
		service.Name("etas"),
		service.Version("latest"),
	)

	// Connect to GoogleMaps
	cf, err := config.Get("google.maps.apikey")
	if err != nil {
		logger.Fatalf("Error loading config: %v", err)
	}
	key := cf.String("")
	if len(key) == 0 {
		logger.Fatalf("Missing require config: google.maps.apikey")
	}
	m, err := maps.NewClient(maps.WithAPIKey(key))
	if err != nil {
		logger.Fatal(err)
	}

	// Register handler
	pb.RegisterETAsHandler(srv.Server(), &handler.ETAs{Maps: m})

	// Run service
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
