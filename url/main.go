package main

import (
	"github.com/micro/services/pkg/tracing"
	"github.com/micro/services/url/handler"
	pb "github.com/micro/services/url/proto"

	"github.com/micro/micro/v3/service"
	"github.com/micro/micro/v3/service/logger"
)

func main() {
	// Create service
	srv := service.New(
		service.Name("url"),
		service.Version("latest"),
	)

	// Register handler
	pb.RegisterUrlHandler(srv.Server(), handler.NewUrl())

	traceCloser := tracing.SetupOpentracing("url")
	defer traceCloser.Close()
	// Run service
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
