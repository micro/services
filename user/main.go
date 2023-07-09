package main

import (
	"micro.dev/v4/service"
	"micro.dev/v4/service/logger"
	"micro.dev/v4/service/store"

	otp "github.com/micro/services/otp/proto"
	adminpb "github.com/micro/services/pkg/service/proto"
	"github.com/micro/services/user/handler"
	proto "github.com/micro/services/user/proto"
)

func main() {
	srv := service.New(
		service.Name("user"),
	)
	srv.Init()

	hd := handler.NewUser(
		store.DefaultStore,
		otp.NewOtpService("otp", srv.Client()),
	)

	proto.RegisterUserHandler(srv.Server(), hd)
	adminpb.RegisterAdminHandler(srv.Server(), hd)

	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
