package main

import (
	"github.com/micro/services/email/handler"
	pb "github.com/micro/services/email/proto"
	"micro.dev/v4/service"
	"micro.dev/v4/service/logger"
)

func main() {
	// Create service
	srv := service.New(
		service.Name("email"),
		service.Version("latest"),
	)

	// Register handler
	pb.RegisterEmailHandler(srv.Server(), handler.NewEmailHandler(srv))

	// Run service
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
