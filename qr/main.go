package main

import (
	"github.com/micro/services/pkg/tracing"
	"github.com/micro/services/qr/handler"
	pb "github.com/micro/services/qr/proto"

	"github.com/micro/micro/v3/service"
	"github.com/micro/micro/v3/service/logger"
)

func main() {
	// Create service
	srv := service.New(
		service.Name("qr"),
		service.Version("latest"),
	)

	// Register handler
	pb.RegisterQrHandler(srv.Server(), handler.New())
	traceCloser := tracing.SetupOpentracing("qr")
	defer traceCloser.Close()

	// Run service
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
