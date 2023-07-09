package main

import (
	"github.com/micro/services/helloworld/handler"
	pb "github.com/micro/services/helloworld/proto"
	"micro.dev/v4/service"
	"micro.dev/v4/service/logger"
)

func main() {
	// Create service
	srv := service.New(
		service.Name("helloworld"),
	)

	// Register Handler
	pb.RegisterHelloworldHandler(srv.Server(), new(handler.Helloworld))

	// Run the service
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
