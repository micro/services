package main

import (
	"github.com/micro/services/wallet/handler"
	pb "github.com/micro/services/wallet/proto"
	"micro.dev/v4/service"
	"micro.dev/v4/service/logger"
)

func main() {
	// Create service
	srv := service.New(
		service.Name("wallet"),
		service.Version("latest"),
	)

	// Register handler
	pb.RegisterWalletHandler(srv.Server(), handler.NewHandler(srv))

	// Run service
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
