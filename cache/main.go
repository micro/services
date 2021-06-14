package main

import (
	"github.com/micro/services/cache/handler"
	pb "github.com/micro/services/cache/proto"
	"github.com/micro/services/pkg/tracing"

	"github.com/micro/micro/v3/service"
	"github.com/micro/micro/v3/service/logger"
)

func main() {
	// Create service
	srv := service.New(
		service.Name("cache"),
		service.Version("latest"),
	)

	// Register handler
	pb.RegisterCacheHandler(srv.Server(), new(handler.Cache))

	traceCloser := tracing.SetupOpentracing("cache")
	defer traceCloser.Close()
	// Run service
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
