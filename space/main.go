package main

import (
	"github.com/micro/micro/v3/service"
	"github.com/micro/micro/v3/service/api"
	"github.com/micro/micro/v3/service/logger"
	admin "github.com/micro/services/pkg/service/proto"
	"github.com/micro/services/space/handler"
)

func main() {
	// Create service
	srv := service.New(
		service.Name("space"),
		service.Version("latest"),
	)

	h := handler.NewSpace(srv)
	// Register handler
	admin.RegisterAdminHandler(srv.Server(), h)
	srv.Server().Handle(
		srv.Server().NewHandler(
			h,
			api.WithEndpoint(
				&api.Endpoint{
					Name:    "Space.Download",
					Handler: "api",
					Method:  []string{"POST", "GET"},
					Path:    []string{"/space/download"},
				}),
		))

	// Run service
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
