package main

import (
	"github.com/micro/services/pkg/tracing"
	"github.com/micro/services/sentiment/handler"
	pb "github.com/micro/services/sentiment/proto"

	"github.com/micro/micro/v3/service"
	"github.com/micro/micro/v3/service/logger"
)

func main() {
	// Create service
	srv := service.New(
		service.Name("sentiment"),
	)

	// Register handler
	pb.RegisterSentimentHandler(srv.Server(), new(handler.Sentiment))
	traceCloser := tracing.SetupOpentracing("sentiment")
	defer traceCloser.Close()

	// Run service
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
