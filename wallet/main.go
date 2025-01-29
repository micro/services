package main

import (
	"github.com/micro/micro/v5/service"
	"github.com/micro/micro/v5/service/logger"
	"github.com/micro/services/wallet/handler"
	pb "github.com/micro/services/wallet/proto"
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
