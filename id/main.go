package main

import (
	"math/rand"
	"time"

	"github.com/micro/services/id/handler"
	pb "github.com/micro/services/id/proto"
	"micro.dev/v4/service"
	"micro.dev/v4/service/logger"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	// Create service
	srv := service.New(
		service.Name("id"),
		service.Version("latest"),
	)

	// Register handler
	pb.RegisterIdHandler(srv.Server(), handler.New())

	// Run service
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
