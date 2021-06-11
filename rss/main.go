package main

import (
	"github.com/micro/micro/v3/service"
	"github.com/micro/micro/v3/service/logger"
	"github.com/micro/services/pkg/tracing"
	"github.com/micro/services/rss/handler"
	pb "github.com/micro/services/rss/proto"
)

func main() {
	// Create service
	srv := service.New(
		service.Name("rss"),
	)

	// Register handler
	pb.RegisterRssHandler(srv.Server(), handler.NewRss())
	traceCloser := tracing.SetupOpentracing("rss")
	defer traceCloser.Close()

	// Run service
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
