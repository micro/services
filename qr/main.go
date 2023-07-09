package main

import (
	admin "github.com/micro/services/pkg/service/proto"
	"github.com/micro/services/qr/handler"
	pb "github.com/micro/services/qr/proto"

	"micro.dev/v4/service"
	"micro.dev/v4/service/logger"
)

func main() {
	// Create service
	srv := service.New(
		service.Name("qr"),
		service.Version("latest"),
	)

	h := handler.New()
	// Register handler
	pb.RegisterQrHandler(srv.Server(), h)
	admin.RegisterAdminHandler(srv.Server(), h)

	// Run service
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
