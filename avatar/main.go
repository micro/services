package main

import (
	"github.com/micro/services/avatar/handler"
	pb "github.com/micro/services/avatar/proto"
	imagePb "github.com/micro/services/image/proto"

	"micro.dev/v4/service"
	"micro.dev/v4/service/logger"
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
