package main

import (
	"github.com/micro/services/sms/handler"
	pb "github.com/micro/services/sms/proto"

	"github.com/micro/micro/v3/service"
	"github.com/micro/micro/v3/service/logger"
)

func main() {
	// Create service
	srv := service.New(
		service.Name("sms"),
		service.Version("latest"),
	)

	// Register handler
	pb.RegisterSmsHandler(srv.Server(), new(handler.Sms))

	// Run service
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
