package main

import (
	"github.com/micro/services/evchargers/handler"
	pb "github.com/micro/services/evchargers/proto"

	"github.com/micro/micro/v5/service"
	"github.com/micro/micro/v5/service/logger"
)

func main() {
	// Create service
	srv := service.New(
		service.Name("evchargers"),
		service.Version("latest"),
	)

	// Register handler
	pb.RegisterEvchargersHandler(srv.Server(), handler.New())

	// Run service
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
