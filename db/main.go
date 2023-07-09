package main

import (
	"database/sql"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/micro/services/db/handler"
	pb "github.com/micro/services/db/proto"
	admin "github.com/micro/services/pkg/service/proto"
	"micro.dev/v4/service"
	"micro.dev/v4/service/config"
	"micro.dev/v4/service/logger"
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

	// Run service
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
