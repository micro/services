package main

import (
	admin "github.com/micro/services/pkg/service/proto"
	"micro.dev/v4/service/store"

	"github.com/micro/services/contact/domain"
	"github.com/micro/services/contact/handler"
	pb "github.com/micro/services/contact/proto"

	"micro.dev/v4/service"
	"micro.dev/v4/service/logger"
)

func main() {
	// Create service
	srv := service.New(
		service.Name("contact"),
		service.Version("latest"),
	)

	contactDomain := domain.NewContactDomain(store.DefaultStore)

	h := handler.NewContact(contactDomain)
	// Register handler
	pb.RegisterContactHandler(srv.Server(), h)
	admin.RegisterAdminHandler(srv.Server(), h)
	// Run service
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
