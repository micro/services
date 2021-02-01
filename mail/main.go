package main

import (
	"github.com/micro/micro/v3/service"
	"github.com/micro/micro/v3/service/logger"

	"github.com/micro/services/mail/handler"
	pb "github.com/micro/services/mail/proto"
)

func main() {
	// Create the service
	srv := service.New(
		service.Name("mail"),
		service.Version("latest"),
	)

	// Register the handler against the server
	pb.RegisterMailHandler(srv.Server(), new(handler.Mail))

	// Run the service
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
