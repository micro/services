package main

import (
	"github.com/micro/micro/v5/service"
	"github.com/micro/micro/v5/service/logger"
	"github.com/micro/services/crypto/handler"
	pb "github.com/micro/services/crypto/proto"
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
