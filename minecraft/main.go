package main

import (
	"github.com/micro/services/minecraft/handler"
	pb "github.com/micro/services/minecraft/proto"

	"github.com/micro/micro/v3/service"
	"github.com/micro/micro/v3/service/logger"
)

func main() {
	// Create service
	srv := service.New(
		service.Name("minecraft"),
		service.Version("latest"),
	)

	// Register handler
	pb.RegisterMinecraftHandler(srv.Server(), new(handler.Minecraft))

	// Run service
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
