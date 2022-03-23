package main

import (
	"github.com/micro/micro/v3/service/api"
	"github.com/micro/services/github/handler"

	"github.com/micro/micro/v3/service"
	"github.com/micro/micro/v3/service/logger"
)

func main() {
	// Create service
	srv := service.New(
		service.Name("github"),
		service.Version("latest"),
	)

	srv.Server().Handle(
		srv.Server().NewHandler(
			handler.NewHandler(srv),
			api.WithEndpoint(
				&api.Endpoint{
					Name:    "Github.Webhook",
					Handler: "api",
					Method:  []string{"POST"},
					Path:    []string{"/github/webhook"},
				}),
		))

	// Run service
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
