package main

import (
	"github.com/micro/micro/v3/service"
	"github.com/micro/micro/v3/service/logger"
	"github.com/micro/services/helloworld/handler"
	pb "github.com/micro/services/helloworld/proto"
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
