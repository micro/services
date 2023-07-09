package main

import (
	"github.com/micro/services/github/handler"
	admin "github.com/micro/services/pkg/service/proto"
	"micro.dev/v4/service/api"

	"micro.dev/v4/service"
	"micro.dev/v4/service/logger"
)

func main() {
	// Create service
	srv := service.New(
		service.Name("github"),
		service.Version("latest"),
	)

	h := handler.NewHandler(srv)
	srv.Server().Handle(
		srv.Server().NewHandler(
			h,
			api.WithEndpoint(
				&api.Endpoint{
					Name:    "Github.Webhook",
					Handler: "api",
					Method:  []string{"POST"},
					Path:    []string{"/github/webhook"},
				}),
		))

	admin.RegisterAdminHandler(srv.Server(), h)
	// Run service
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
