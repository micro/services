package main

import (
	"github.com/micro/services/translation/handler"
	pb "github.com/micro/services/translation/proto"

	"github.com/micro/micro/v3/service"
	"github.com/micro/micro/v3/service/logger"
)

func main() {
	// Create service
	srv := service.New(
		service.Name("translation"),
		service.Version("latest"),
	)

	// Register handler
	pb.RegisterTranslationHandler(srv.Server(), handler.NewTranslation())

	// Run service
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
