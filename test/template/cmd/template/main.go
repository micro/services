package main

import (
	"github.com/micro/services/test/template/handler"

	"micro.dev/v4/service"
	"micro.dev/v4/service/logger"
)

func main() {
	// Create service
	srv := service.New(
		service.Name("idiomatic"),
		service.Version("latest"),
	)

	// Register handler
	srv.Handle(new(handler.Idiomatic))

	// Run service
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
