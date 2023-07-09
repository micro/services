package main

import (
	"github.com/micro/services/sentiment/handler"
	pb "github.com/micro/services/sentiment/proto"

	"micro.dev/v4/service"
	"micro.dev/v4/service/logger"
)

func main() {
	// Create service
	srv := service.New(
		service.Name("sentiment"),
	)

	// Register handler
	pb.RegisterSentimentHandler(srv.Server(), new(handler.Sentiment))

	// Run service
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
