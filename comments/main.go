package main

import (
	"github.com/micro/micro/v3/service"
	log "github.com/micro/micro/v3/service/logger"
	"github.com/micro/services/comments/handler"
	pb "github.com/micro/services/comments/proto"
	admin "github.com/micro/services/pkg/service/proto"
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
