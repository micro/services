package main

import (
	"github.com/micro/services/seen/handler"
	pb "github.com/micro/services/seen/proto"

	"github.com/micro/micro/v3/service"
	"github.com/micro/micro/v3/service/logger"
)

func main() {
	// Create service
	srv := service.New(
		service.Name("seen"),
		service.Version("latest"),
	)

	// Register handler
	pb.RegisterSeenHandler(srv.Server(), new(handler.Seen))

	// Run service
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
