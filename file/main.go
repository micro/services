package main

import (
	"github.com/micro/services/file/handler"
	pb "github.com/micro/services/file/proto"
	admin "github.com/micro/services/pkg/service/proto"

	"github.com/micro/micro/v5/service"
	"github.com/micro/micro/v5/service/logger"
)

func main() {
	// Create service
	srv := service.New(
		service.Name("file"),
		service.Version("latest"),
	)

	h := handler.NewFile()
	// Register handler
	pb.RegisterFileHandler(srv.Server(), h)
	admin.RegisterAdminHandler(srv.Server(), h)

	// Run service
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
