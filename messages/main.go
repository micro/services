package main

import (
	"github.com/micro/micro/v3/service"
	"github.com/micro/micro/v3/service/logger"

	"github.com/micro/services/messages/handler"
	pb "github.com/micro/services/messages/proto"
)

func main() {
	// Create the service
	srv := service.New(
		service.Name("messages"),
		service.Version("latest"),
	)

	// Register the handler against the server
	pb.RegisterMessagesHandler(srv.Server(), new(handler.Messages))

	// Run the service
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
