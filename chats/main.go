package main

import (
	"time"

	"github.com/micro/services/chats/handler"
	pb "github.com/micro/services/chats/proto"

	"github.com/micro/micro/v3/service"
	"github.com/micro/micro/v3/service/logger"
)

func main() {
	// Create service
	srv := service.New(
		service.Name("chats"),
		service.Version("latest"),
	)

	h := &handler.Chats{Time: time.Now}
	// Register handler
	pb.RegisterChatsHandler(srv.Server(), h)

	// Run service
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
