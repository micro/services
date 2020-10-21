package main

import (
	"github.com/micro/micro/service/v3/logger"
	"github.com/micro/micro/service/v3"
	"github.com/micro/services/users/db"
	"github.com/micro/services/users/handler"
	proto "github.com/micro/services/users/proto"
)

func main() {
	service := service.New(
		service.Name("users"),
	)

	service.Init()
	db.Init()

	proto.RegisterUsersHandler(service.Server(), new(handler.Users))

	if err := service.Run(); err != nil {
		logger.Fatal(err)
	}
}
