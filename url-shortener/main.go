package main

import (
	"github.com/micro/services/url-shortener/handler"
	pb "github.com/micro/services/url-shortener/proto"

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
	pb.RegisterUrlShortenerHandler(srv.Server(), handler.NewUrlShortener())

	// Run service
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
