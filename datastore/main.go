package main

import (
	"github.com/micro/services/datastore/handler"
	pb "github.com/micro/services/datastore/proto"

	"github.com/micro/micro/v3/service"
	"github.com/micro/micro/v3/service/logger"
)

func main() {
	// Create service
	srv := service.New(
		service.Name("datastore"),
		service.Version("latest"),
	)

	// Register handler
	pb.RegisterDatastoreHandler(srv.Server(), new(handler.Datastore))

	// Run service
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
