package main

import (
	"github.com/micro/services/comments/handler"
	pb "github.com/micro/services/comments/proto"
	admin "github.com/micro/services/pkg/service/proto"
	"micro.dev/v4/service"
	log "micro.dev/v4/service/logger"
)

func main() {
	// New Service
	srv := service.New(
		service.Name("comments"),
		service.Version("latest"),
	)

	// Initialise service
	srv.Init()

	h := handler.New(srv.Client())
	// Register Handler
	pb.RegisterCommentsHandler(srv.Server(), h)
	admin.RegisterAdminHandler(srv.Server(), h)

	// Run service
	if err := srv.Run(); err != nil {
		log.Fatal(err)
	}
}
