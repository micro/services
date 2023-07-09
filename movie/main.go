package main

import (
	"github.com/micro/services/movie/handler"
	pb "github.com/micro/services/movie/proto"

	"micro.dev/v4/service"
	"micro.dev/v4/service/logger"
)

func main() {
	// Create service
	srv := service.New(
		service.Name("movie"),
		service.Version("latest"),
	)

	// Register handler
	pb.RegisterMovieHandler(srv.Server(), handler.New())

	// Run service
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
