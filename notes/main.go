package main

import (
	"github.com/micro/micro/v3/service"
	"github.com/micro/micro/v3/service/logger"

	"github.com/micro/services/notes/handler"
	pb "github.com/micro/services/notes/proto"
)

func main() {
	srv := service.New(
		service.Name("notes"),
		service.Version("latest"),
	)

	pb.RegisterNotesHandler(srv.Server(), handler.New())

	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
