package main

import (
	"github.com/micro/services/users/handler"
	proto "github.com/micro/services/users/proto"
	"github.com/micro/micro/v3/service"
	"github.com/micro/micro/v3/service/logger"
)

func main() {
	service := service.New(
		service.Name("users"),
	)

	service.Init()

	proto.RegisterUsersHandler(service.Server(), handler.NewUsers())

	if err := service.Run(); err != nil {
		logger.Fatal(err)
	}
}
