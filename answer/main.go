package main

import (
	"github.com/micro/micro/v3/service"
	"github.com/micro/micro/v3/service/logger"
	"github.com/micro/services/answer/handler"
	pb "github.com/micro/services/answer/proto"
)

func main() {
	// Create service
	srv := service.New(
		service.Name("answer"),
		service.Version("latest"),
	)

	// Register handler
	pb.RegisterAnswerHandler(srv.Server(), new(handler.Answer))

	// Run service
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
