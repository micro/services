package main

import (
	"time"

	"github.com/micro/services/stock/handler"
	pb "github.com/micro/services/stock/proto"

	"github.com/patrickmn/go-cache"
	"micro.dev/v4/service"
	"micro.dev/v4/service/config"
	"micro.dev/v4/service/logger"
)

func main() {
	// Create service
	srv := service.New(
		service.Name("stock"),
		service.Version("latest"),
	)

	v, err := config.Get("finage.api")
	if err != nil {
		logger.Fatalf("finage.api config not found: %v", err)
	}
	api := v.String("")
	if len(api) == 0 {
		logger.Fatal("finage.api config not found")
	}
	v, err = config.Get("finage.key")
	if err != nil {
		logger.Fatalf("finage.key config not found: %v", err)
	}
	key := v.String("")
	if len(key) == 0 {
		logger.Fatal("finage.key config not found")
	}

	// Register handler
	pb.RegisterStockHandler(srv.Server(), &handler.Stock{
		Api:   api,
		Key:   key,
		Cache: cache.New(5*time.Minute, 10*time.Minute),
	})

	// Run service
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
