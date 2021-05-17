package main

import (
	"github.com/micro/services/url/handler"
	pb "github.com/micro/services/url/proto"

	"github.com/micro/micro/v3/service"
	"github.com/micro/micro/v3/service/logger"
)

func main() {
	// Create service
	srv := service.New(
		service.Name("url-shortener"),
		service.Version("latest"),
	)

	// Register handler
	pb.RegisterUrlShortenerHandler(srv.Server(), handler.NewUrl())

	// Run service
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
