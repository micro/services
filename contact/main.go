package main

import (
	"github.com/micro/micro/v3/service/store"

	"github.com/micro/services/contact/domain"
	"github.com/micro/services/contact/handler"
	pb "github.com/micro/services/contact/proto"

	"github.com/micro/micro/v3/service"
	"github.com/micro/micro/v3/service/logger"
)

func main() {
	// Create service
	srv := service.New(
		service.Name("contact"),
		service.Version("latest"),
	)

	contactDomain := domain.NewContactDomain(store.DefaultStore)

	// Register handler
	pb.RegisterContactHandler(srv.Server(), handler.NewContact(contactDomain))

	// Run service
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
