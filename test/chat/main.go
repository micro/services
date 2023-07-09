package main

import (
	"micro.dev/v4/service"
	"micro.dev/v4/service/logger"

	"github.com/micro/services/test/chat/handler"
	pb "github.com/micro/services/test/chat/proto"
)

func main() {
	// Create the service
	srv := service.New(
		service.Name("chat"),
		service.Version("latest"),
	)

	// Register the handler against the server
	pb.RegisterChatHandler(srv.Server(), new(handler.Chat))

	// Run the service
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
