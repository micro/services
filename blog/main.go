package main

import (
	"github.com/micro/services/blog/handler"
	pb "github.com/micro/services/blog/proto"

	"github.com/micro/micro/v3/service"
	"github.com/micro/micro/v3/service/logger"
)

func main() {
	// Create service
	srv := service.New(
		service.Name("blog"),
	)

	// Register handler
	pb.RegisterBlogHandler(srv.Server(), new(handler.Blog))

	// Run service
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
