package main

import (
	"time"

	"github.com/micro/micro/v3/service"
	"github.com/micro/micro/v3/service/config"
	"github.com/micro/micro/v3/service/logger"
	"github.com/micro/micro/v3/service/store"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	authPb "github.com/micro/micro/v3/proto/auth"

	db "github.com/micro/services/db/proto"
	otp "github.com/micro/services/otp/proto"
	"github.com/micro/services/pkg/tracing"
	"github.com/micro/services/user/handler"
	"github.com/micro/services/user/migrate"
	proto "github.com/micro/services/user/proto"
)

var pgxDsn = "postgresql://postgres:postgres@localhost:5432/db?sslmode=disable"

func migrateData(from db.DbService, to store.Store, authAccount authPb.AccountsService) {
	startTime := time.Now()
	logger.Info("start migrate ...")
	defer func() {
		logger.Infof("migrate finish, use time: %v", time.Since(startTime))
	}()

	// Connect to the database
	cfg, err := config.Get("micro.db.database")
	if err != nil {
		logger.Fatalf("Error loading config: %v", err)
	}

	dsn := cfg.String(pgxDsn)
	gormDb, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		logger.Fatal("Failed to connect to ")
	}

	migration := migrate.NewMigration(gormDb)
	migration.Do()

	return
}

func main() {
	srv := service.New(
		service.Name("user"),
	)
	srv.Init()

	from := db.NewDbService("db", srv.Client())
	authAccount := authPb.NewAccountsService("auth", srv.Client())
	go migrateData(from, store.DefaultStore, authAccount)

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
