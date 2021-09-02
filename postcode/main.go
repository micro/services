package main

import (
	"github.com/micro/micro/v3/service"
	"github.com/micro/micro/v3/service/logger"
	"github.com/micro/services/postcode/handler"
	pb "github.com/micro/services/postcode/proto"
)

func main() {
	// Create service
	srv := service.New(
		service.Name("postcode"),
		service.Version("latest"),
	)

	// Register handler
	pb.RegisterPostcodeHandler(srv.Server(), new(handler.Postcode))

	// Run service
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
