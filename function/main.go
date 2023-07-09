package main

import (
	"github.com/micro/services/function/handler"
	pb "github.com/micro/services/function/proto"
	admin "github.com/micro/services/pkg/service/proto"
	"micro.dev/v4/service"
	"micro.dev/v4/service/logger"
)

func main() {
	// Create service
	srv := service.New(
		service.Name("function"),
		service.Version("latest"),
	)

	h := handler.NewFunction(srv)
	// Register handler
	pb.RegisterFunctionHandler(srv.Server(), h)
	admin.RegisterAdminHandler(srv.Server(), h)

	// Run service
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
