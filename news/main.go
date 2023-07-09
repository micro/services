package main

import (
	"github.com/micro/services/news/handler"
	pb "github.com/micro/services/news/proto"

	"micro.dev/v4/service"
	"micro.dev/v4/service/config"
	"micro.dev/v4/service/logger"
)

func main() {
	// Create service
	srv := service.New(
		service.Name("news"),
		service.Version("latest"),
	)

	// Setup google maps
	c, err := config.Get("news.apikey")
	if err != nil {
		logger.Fatalf("Error loading config: %v", err)
	}

	apiKey := c.String("")
	if len(apiKey) == 0 {
		logger.Fatalf("Missing required config: news.apikey")
	}

	// Register handler
	pb.RegisterNewsHandler(srv.Server(), handler.New(apiKey))

	// Run service
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
