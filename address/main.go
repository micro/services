package main

import (
	"github.com/micro/services/address/handler"
	pb "github.com/micro/services/address/proto"
	"micro.dev/v4/service"
	"micro.dev/v4/service/config"
	"micro.dev/v4/service/logger"
)

func main() {
	// Create service
	srv := service.New(
		service.Name("address"),
		service.Version("latest"),
	)

	v, err := config.Get("address.api")
	if err != nil {
		logger.Fatalf("address.api config not found: %v", err)
	}
	api := v.String("")
	if len(api) == 0 {
		logger.Fatal("address.api config not found")
	}
	v, err = config.Get("address.key")
	if err != nil {
		logger.Fatalf("address.key config not found: %v", err)
	}
	key := v.String("")
	if len(key) == 0 {
		logger.Fatal("address.key config not found")
	}

	// Register handler
	pb.RegisterAddressHandler(srv.Server(), &handler.Address{
		Url: api,
		Key: key,
	})

	// Run service
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
