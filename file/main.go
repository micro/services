package main

import (
	"github.com/micro/services/file/handler"
	pb "github.com/micro/services/file/proto"
	"github.com/micro/services/pkg/tracing"

	"github.com/micro/micro/v3/service"
	"github.com/micro/micro/v3/service/logger"
)

func main() {
	// Create service
	srv := service.New(
		service.Name("file"),
		service.Version("latest"),
	)

	// Register handler
	pb.RegisterFileHandler(srv.Server(), handler.NewFile())

	traceCloser := tracing.SetupOpentracing("file")
	defer traceCloser.Close()

	// Run service
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
