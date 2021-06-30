package main

import (
	"github.com/micro/micro/v3/service"
	"github.com/micro/micro/v3/service/logger"
	"github.com/micro/services/ask/handler"
	pb "github.com/micro/services/ask/proto"
)

func main() {
	// Create service
	srv := service.New(
		service.Name("ask"),
		service.Version("latest"),
	)

	// Register handler
	pb.RegisterAskHandler(srv.Server(), new(handler.Ask))

	// Run service
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
