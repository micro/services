package main

import (
	"github.com/micro/micro/v3/service"
	"github.com/micro/micro/v3/service/logger"
	"github.com/micro/services/groups/handler"
	pb "github.com/micro/services/groups/proto"
)

func main() {
	// Create service
	srv := service.New(
		service.Name("groups"),
		service.Version("latest"),
	)

	// Register handler
	pb.RegisterGroupsHandler(srv.Server(), new(handler.Groups))

	// Run service
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
