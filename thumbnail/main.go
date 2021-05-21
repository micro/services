package main

import (
	iproto "github.com/micro/services/image/proto"
	"github.com/micro/services/thumbnail/handler"
	pb "github.com/micro/services/thumbnail/proto"

	"github.com/micro/micro/v3/service"
	"github.com/micro/micro/v3/service/logger"
)

func main() {
	// Create service
	srv := service.New(
		service.Name("thumbnail"),
		service.Version("latest"),
	)

	// Register handler
	pb.RegisterThumbnailHandler(srv.Server(), handler.NewThumbnail(iproto.NewImageService("image", srv.Client())))

	// Run service
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
