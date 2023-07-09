package main

import (
	admin "github.com/micro/services/pkg/service/proto"
	"github.com/micro/services/space/handler"
	"micro.dev/v4/service"
	"micro.dev/v4/service/api"
	"micro.dev/v4/service/logger"
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
