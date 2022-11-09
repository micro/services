package main

import (
	"database/sql"

	pb "github.com/micro/services/db/proto"
	admin "github.com/micro/services/pkg/service/proto"
	"github.com/micro/services/pkg/tracing"
	"github.com/micro/services/db/handler"
	"github.com/micro/micro/v3/service"
	"github.com/micro/micro/v3/service/logger"
	"github.com/micro/micro/v3/service/config"
	_ "github.com/jackc/pgx/v4/stdlib"
)

var dbAddress = "postgresql://postgres:postgres@localhost:5432/db?sslmode=disable"

func main() {
	// Create service
	srv := service.New(
		service.Name("db"),
		service.Version("latest"),
	)

	// Connect to the database
	cfg, err := config.Get("db.address")
	if err != nil {
		logger.Fatalf("Error loading config: %v", err)
	}
	addr := cfg.String(dbAddress)
	sqlDB, err := sql.Open("pgx", addr)
	if err != nil {
		logger.Fatalf("Failed to open connection to DB %s", err)
	}
	h := &handler.Db{}
	h.DBConn(sqlDB)

	// Register handler
	pb.RegisterDbHandler(srv.Server(), h)
	admin.RegisterAdminHandler(srv.Server(), h)

	traceCloser := tracing.SetupOpentracing("db")
	defer traceCloser.Close()
	// Run service
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
