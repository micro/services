package main

import (
	"time"

	"github.com/micro/micro/v3/service"
	"github.com/micro/micro/v3/service/logger"
	"github.com/micro/micro/v3/service/store"

	db "github.com/micro/services/db/proto"
	otp "github.com/micro/services/otp/proto"
	"github.com/micro/services/pkg/tracing"
	"github.com/micro/services/user/handler"
	"github.com/micro/services/user/migrate"
	proto "github.com/micro/services/user/proto"
)

func migrateData(from db.DbService, to store.Store) {
	startTime := time.Now()
	logger.Info("start migrate ...")
	defer func() {
		logger.Infof("migrate finish, use time: %v", time.Since(startTime))
	}()

	// users
	u := migrate.NewUserMigration(from, to)
	err := u.Do()
	if err != nil {
		logger.Errorf("migrate users data error: %v", err)
	}

	//

	return
}

func main() {
	srv := service.New(
		service.Name("user"),
	)
	srv.Init()

	from := db.NewDbService("db", srv.Client())
	go migrateData(from, store.DefaultStore)

	hd := handler.NewUser(
		store.DefaultStore,
		otp.NewOtpService("otp", srv.Client()),
	)

	proto.RegisterUserHandler(srv.Server(), hd)
	traceCloser := tracing.SetupOpentracing("user")
	defer traceCloser.Close()

	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
