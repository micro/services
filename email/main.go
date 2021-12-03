package main

import (
	"github.com/micro/micro/v3/service"
	"github.com/micro/micro/v3/service/logger"
	"github.com/micro/services/email/handler"
	pb "github.com/micro/services/email/proto"
	"github.com/micro/services/pkg/tracing"
)

func main() {
	// Create service
	srv := service.New(
		service.Name("email"),
		service.Version("latest"),
	)

	// Register handler
	pb.RegisterEmailHandler(srv.Server(), handler.NewEmailHandler(srv))
	traceCloser := tracing.SetupOpentracing("email")
	defer traceCloser.Close()

	// Run service
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
