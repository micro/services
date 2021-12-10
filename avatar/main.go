package main

import (
	"github.com/micro/services/avatar/handler"
	pb "github.com/micro/services/avatar/proto"

	"github.com/micro/micro/v3/service"
	"github.com/micro/micro/v3/service/logger"
)

func main() {
	// Create service
	srv := service.New(
		service.Name("avatar"),
		service.Version("latest"),
	)

	// Register handler
	pb.RegisterAvatarHandler(srv.Server(), new(handler.Avatar))

	// Run service
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
