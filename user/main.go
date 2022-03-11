package main

import (
	"github.com/micro/micro/v3/service"
	"github.com/micro/micro/v3/service/logger"
	"github.com/micro/micro/v3/service/store"

	otp "github.com/micro/services/otp/proto"
	adminpb "github.com/micro/services/pkg/service/proto"
	"github.com/micro/services/pkg/tracing"
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
	proto.RegisterUserAdminHandler(srv.Server(), hd)
	traceCloser := tracing.SetupOpentracing("user")
	defer traceCloser.Close()

	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
