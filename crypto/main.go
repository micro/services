package main

import (
	"github.com/micro/services/crypto/handler"
	pb "github.com/micro/services/crypto/proto"
	"micro.dev/v4/service"
	"micro.dev/v4/service/logger"
)

func main() {
	// Create service
	srv := service.New(
		service.Name("crypto"),
		service.Version("latest"),
	)

	// Register handler
	pb.RegisterCryptoHandler(srv.Server(), handler.New())

	// Run service
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
