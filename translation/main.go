package main

import (
	"github.com/micro/micro/v3/service/config"

	"github.com/micro/services/translation/domain"
	"github.com/micro/services/translation/handler"
	pb "github.com/micro/services/translation/proto"

	"github.com/micro/micro/v3/service"
	"github.com/micro/micro/v3/service/logger"
)

func youdaoConfig() domain.YoudaoConfig {
	v, err := config.Get("translation.youdao")
	if err != nil {
		logger.Fatalf("translation.youdao config not found: %v", err)
	}
	cfg := v.StringMap(map[string]string{})

	youdaoCfg := domain.YoudaoConfig{}
	youdaoCfg.Api = cfg["api"]
	youdaoCfg.AppKey = cfg["appKey"]
	youdaoCfg.Secret = cfg["secret"]
	if youdaoCfg.Api == "" || youdaoCfg.AppKey == "" || youdaoCfg.Secret == "" {
		logger.Fatalf("translation.youdao configurations miss some fields: %+v", youdaoCfg)
	}

	return youdaoCfg
}

func main() {
	// Create service
	srv := service.New(
		service.Name("translation"),
		service.Version("latest"),
	)

	// Register handler
	pb.RegisterTranslationHandler(srv.Server(), handler.NewTranslation(youdaoConfig()))

	// Run service
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
