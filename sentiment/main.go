package main

import (
	"github.com/micro/services/sentiment/handler"
	pb "github.com/micro/services/sentiment/proto"

	"github.com/cdipaolo/sentiment"
	"github.com/micro/micro/v3/service"
	"github.com/micro/micro/v3/service/logger"
)

func main() {
	// Create service
	srv := service.New(
		service.Name("sentiment"),
	)

	// load sentiment analysis tool
	md, err := sentiment.Restore()
	if err != nil {
		logger.Fatal(err)
	}

	// Register handler
	pb.RegisterSentimentHandler(srv.Server(), &handler.Sentiment{&md})

	// Run service
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
