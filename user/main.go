package main

import (
	"time"

	"github.com/micro/micro/v3/service"
	"github.com/micro/micro/v3/service/config"
	"github.com/micro/micro/v3/service/logger"
	"github.com/micro/micro/v3/service/store"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	otp "github.com/micro/services/otp/proto"
	adminpb "github.com/micro/services/pkg/service/proto"
	"github.com/micro/services/pkg/tracing"
	"github.com/micro/services/user/handler"
	"github.com/micro/services/user/migrate"
	proto "github.com/micro/services/user/proto"
)

var pgxDsn = "postgresql://postgres:postgres@localhost:5432/db?sslmode=disable"

func migrateData() {
	startTime := time.Now()
	logger.Info("start migrate ...")
	defer func() {
		logger.Infof("all migrations are finished, use time: %v", time.Since(startTime))
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
	if err := migration.Do(); err != nil {
		logger.Fatal("migrate error: ", err)
	}

	return
}

func main() {
	srv := service.New(
		service.Name("user"),
	)
	srv.Init()

	// migration work
	go migrateData()

	hd := handler.NewUser(
		store.DefaultStore,
		otp.NewOtpService("otp", srv.Client()),
	)

	proto.RegisterUserHandler(srv.Server(), hd)
	adminpb.RegisterAdminHandler(srv.Server(), hd)
	traceCloser := tracing.SetupOpentracing("user")
	defer traceCloser.Close()

	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
