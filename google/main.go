package main

import (
	"github.com/micro/micro/v3/service"
	"github.com/micro/micro/v3/service/config"
	"github.com/micro/micro/v3/service/logger"
	"github.com/micro/services/google/handler"
	pb "github.com/micro/services/google/proto"
)

func main() {
	// Create service
	srv := service.New(
		service.Name("google"),
		service.Version("latest"),
	)

	// Setup google maps
	c, err := config.Get("google.apikey")
	if err != nil {
		logger.Fatalf("Error loading config: %v", err)
	}
	apiKey := c.String("")
	if len(apiKey) == 0 {
		logger.Fatalf("Missing required config: google.apikey")
	}
	// Setup google maps
	c, err = config.Get("google.cx_id")
	if err != nil {
		logger.Fatalf("Error loading config: %v", err)
	}
	cxId := c.String("")
	if len(cxId) == 0 {
		logger.Fatalf("Missing required config: google.cxId")
	}

	// Register handler
	pb.RegisterGoogleHandler(srv.Server(), handler.New(apiKey, cxId))

	// Run service
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
