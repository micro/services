package main

import (
	"github.com/micro/services/function/handler"
	pb "github.com/micro/services/function/proto"

	"github.com/micro/micro/v3/service"
	"github.com/micro/micro/v3/service/logger"
)

func main() {
	// Create service
	srv := service.New(
		service.Name("function"),
		service.Version("latest"),
	)

	// Register handler
	pb.RegisterFunctionHandler(srv.Server(), new(handler.Function))

	// Run service
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
