package main

import (
	"github.com/micro/services/otp/handler"
	pb "github.com/micro/services/otp/proto"
	admin "github.com/micro/services/pkg/service/proto"
	"github.com/micro/services/pkg/tracing"

	"github.com/micro/micro/v3/service"
	"github.com/micro/micro/v3/service/logger"
)

func main() {
	// Create service
	srv := service.New(
		service.Name("otp"),
	)

	h := new(handler.Otp)
	// Register handler
	pb.RegisterOtpHandler(srv.Server(), h)
	admin.RegisterAdminHandler(srv.Server(), h)

	traceCloser := tracing.SetupOpentracing("otp")
	defer traceCloser.Close()
	// Run service
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
