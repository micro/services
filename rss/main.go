package main

import (
	pb "github.com/micro/services/rss/proto"
	"github.com/micro/services/rss/handler"
	"github.com/micro/micro/v3/service"
	"github.com/micro/micro/v3/service/logger"
)

func main() {
	// Create service
	srv := service.New(
		service.Name("rss"),
	)

	// Register handler
	pb.RegisterRssHandler(srv.Server(), handler.NewRss())

	// Run service
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
