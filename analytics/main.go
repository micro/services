package main

import (
	"analytics/handler"
	pb "analytics/proto"

	"github.com/micro/micro/v3/service"
	"github.com/micro/micro/v3/service/logger"
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
