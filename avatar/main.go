package main

import (
	"github.com/micro/services/avatar/handler"
	pb "github.com/micro/services/avatar/proto"
	imagePb "github.com/micro/services/image/proto"

	"github.com/micro/micro/v5/service"
	"github.com/micro/micro/v5/service/logger"
)

func main() {
	// Create service
	srv := service.New(
		service.Name("avatar"),
		service.Version("latest"),
	)

	// Register handler
	hdlr := handler.NewAvatar(imagePb.NewImageService("image", srv.Client()))
	pb.RegisterAvatarHandler(srv.Server(), hdlr)

	// Run service
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
