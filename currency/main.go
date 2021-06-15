package main

import (
	"time"

	"github.com/micro/micro/v3/service"
	"github.com/micro/micro/v3/service/config"
	"github.com/micro/micro/v3/service/logger"
	"github.com/micro/services/currency/handler"
	pb "github.com/micro/services/currency/proto"
	"github.com/patrickmn/go-cache"
)

func main() {
	// Create service
	srv := service.New(
		service.Name("currency"),
		service.Version("latest"),
	)

	v, err := config.Get("exchangerate.api")
	if err != nil {
		logger.Fatalf("exchangerate.api config not found: %v", err)
	}
	api := v.String("")
	if len(api) == 0 {
		logger.Fatal("exchangerate.api config not found")
	}
	v, err = config.Get("exchangerate.key")
	if err != nil {
		logger.Fatalf("exchangerate.key config not found: %v", err)
	}
	key := v.String("")
	if len(key) == 0 {
		logger.Fatal("exchangerate.key config not found")
	}

	// Register handler
	pb.RegisterCurrencyHandler(srv.Server(), &handler.Currency{
		Api:   api + key,
		Cache: cache.New(5*time.Minute, 10*time.Minute),
	})

	// Run service
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
