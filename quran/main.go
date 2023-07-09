package main

import (
	"github.com/micro/services/quran/handler"
	pb "github.com/micro/services/quran/proto"

	"micro.dev/v4/service"
	"micro.dev/v4/service/logger"
)

func main() {
	// Create service
	srv := service.New(
		service.Name("quran"),
		service.Version("latest"),
	)

	// Register handler
	pb.RegisterQuranHandler(srv.Server(), handler.New())

	// Run service
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
