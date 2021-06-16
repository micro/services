package main

import (
	"github.com/micro/micro/v3/service"
	"github.com/micro/micro/v3/service/api"
	"github.com/micro/micro/v3/service/logger"
	db "github.com/micro/services/db/proto"
	"github.com/micro/services/pkg/tracing"
	"github.com/micro/services/user/handler"
	proto "github.com/micro/services/user/proto"
)

func main() {
	service := service.New(
		service.Name("user1"),
	)

	service.Init()

	handl := handler.NewUser(db.NewDbService("db", service.Client()))
	service.Server().Handle(
		service.Server().NewHandler(
			handl,
			api.WithEndpoint(
				&api.Endpoint{
					Name:    "Verify",
					Handler: "rpc",
					Method:  []string{"GET", "POST", "OPTIONS", "PUT", "HEAD", "DELETE"},
					Path:    []string{"^/v1/.*$"},
					Stream:  true,
				}),
		))

	proto.RegisterUserHandler(service.Server(), handl)
	traceCloser := tracing.SetupOpentracing("user")
	defer traceCloser.Close()

	if err := service.Run(); err != nil {
		logger.Fatal(err)
	}
}
