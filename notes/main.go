package main

import (
	"github.com/micro/micro/v5/service"
	log "github.com/micro/micro/v5/service/logger"
	"github.com/micro/services/notes/handler"
	pb "github.com/micro/services/notes/proto"
	admin "github.com/micro/services/pkg/service/proto"
)

func main() {
	// New Service
	srv := service.New(
		service.Name("notes"),
		service.Version("latest"),
	)

	// Initialise service
	srv.Init()

	h := handler.New(srv.Client())
	// Register Handler
	pb.RegisterNotesHandler(srv.Server(), h)
	admin.RegisterAdminHandler(srv.Server(), h)

	// Run service
	if err := srv.Run(); err != nil {
		log.Fatal(err)
	}
}
