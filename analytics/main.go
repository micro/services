package main

import (
	"github.com/micro/micro/v3/service"
	"github.com/micro/micro/v3/service/logger"
	"github.com/micro/services/analytics/handler"
	pb "github.com/micro/services/analytics/proto"
)

func main() {
	// Create service
	srv := service.New(
		service.Name("analytics"),
		service.Version("latest"),
	)

	h := handler.New()

	// Register handler
	pb.RegisterAnalyticsHandler(srv.Server(), h)

	// Run service
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
