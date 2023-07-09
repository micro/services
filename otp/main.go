package main

import (
	"github.com/micro/services/otp/handler"
	pb "github.com/micro/services/otp/proto"
	admin "github.com/micro/services/pkg/service/proto"

	"micro.dev/v4/service"
	"micro.dev/v4/service/logger"
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

	// Run service
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
