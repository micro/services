package main

import (
	"github.com/micro/services/youtube/handler"
	pb "github.com/micro/services/youtube/proto"
	"micro.dev/v4/service"
	"micro.dev/v4/service/config"
	"micro.dev/v4/service/logger"
)

func main() {
	// Create service
	srv := service.New(
		service.Name("youtube"),
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

	// Register handler
	pb.RegisterYoutubeHandler(srv.Server(), handler.New(apiKey))

	// Run service
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
