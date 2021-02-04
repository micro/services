package main

import (
	"github.com/micro/services/seen/domain"
	"github.com/micro/services/seen/handler"
	pb "github.com/micro/services/seen/proto"

	"github.com/micro/micro/v3/service"
	"github.com/micro/micro/v3/service/logger"
	"github.com/micro/micro/v3/service/store"
)

func main() {
	// Create service
	srv := service.New(
		service.Name("seen"),
		service.Version("latest"),
	)

	// Register handler
	pb.RegisterSeenHandler(srv.Server(), &handler.Seen{
		Domain: domain.New(store.DefaultStore),
	})

	// Run service
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
