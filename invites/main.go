package main

import (
	"github.com/micro/services/invites/handler"
	pb "github.com/micro/services/invites/proto"

	"github.com/micro/micro/v3/service"
	"github.com/micro/micro/v3/service/logger"
)

func main() {
	// Create service
	srv := service.New(
		service.Name("invites"),
		service.Version("latest"),
	)

	// Register handler
	pb.RegisterInvitesHandler(srv.Server(), new(handler.Invites))

	// Run service
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
