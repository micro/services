package main

import (
	"github.com/micro/services/dns/handler"
	pb "github.com/micro/services/dns/proto"
	"micro.dev/v4/service"
	"micro.dev/v4/service/logger"
)

func main() {
	// Create service
	srv := service.New(
		service.Name("dns"),
	)

	// Register handler
	pb.RegisterDnsHandler(srv.Server(), handler.New())

	// Run service
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
