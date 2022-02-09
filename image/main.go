package main

import (
	"github.com/micro/services/image/handler"
	pb "github.com/micro/services/image/proto"
	admin "github.com/micro/services/pkg/service/proto"
	"github.com/micro/services/pkg/tracing"

	"github.com/micro/micro/v3/service"
	"github.com/micro/micro/v3/service/logger"
)

func main() {
	// Create service
	srv := service.New(
		service.Name("image"),
		service.Version("latest"),
	)

	h := handler.NewImage()
	// Register handler
	pb.RegisterImageHandler(srv.Server(), h)
	admin.RegisterAdminHandler(srv.Server(), h)
	traceCloser := tracing.SetupOpentracing("image")
	defer traceCloser.Close()

	// Run service
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
