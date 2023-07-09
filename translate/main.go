package main

import (
	"github.com/micro/services/translate/handler"
	pb "github.com/micro/services/translate/proto"

	"micro.dev/v4/service"
	"micro.dev/v4/service/logger"
)

func main() {
	// Create service
	srv := service.New(
		service.Name("translate"),
		service.Version("latest"),
	)

	// Register handler
	pb.RegisterTranslateHandler(srv.Server(), handler.NewTranslation())

	// Run service
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
