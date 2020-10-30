package main

import (
	"blog/handler"
	pb "blog/proto"

	"github.com/micro/micro/v3/service"
	"github.com/micro/micro/v3/service/logger"
)

func main() {
	// Create service
	srv := service.New(
		service.Name("blog"),
		service.Version("latest"),
	)

	// Register handler
	pb.RegisterBlogHandler(srv.Server(), new(handler.Blog))

	// Run service
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
