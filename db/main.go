package main

import (
	pb "db/proto"

	"github.com/micro/services/db/handler"

	"github.com/micro/micro/v3/service"
	"github.com/micro/micro/v3/service/logger"

	"database/sql"

	"github.com/micro/micro/v3/service/config"

	_ "github.com/jackc/pgx/v4/stdlib"
)

var dbAddress = "postgresql://postgres:postgres@localhost:5432/invites?sslmode=disable"

func main() {
	// Create service
	srv := service.New(
		service.Name("db"),
		service.Version("latest"),
	)

	// Connect to the database
	cfg, err := config.Get("invites.database")
	if err != nil {
		logger.Fatalf("Error loading config: %v", err)
	}
	addr := cfg.String(dbAddress)
	sqlDB, err := sql.Open("pgx", addr)
	if err != nil {
		logger.Fatalf("Failed to open connection to DB %s", err)
	}
	h := &handler.Db{}
	h.DBConn(sqlDB).Migrations(&handler.Db{})

	// Register handler
	pb.RegisterInvitesHandler(srv.Server(), h)

	// Register handler
	pb.RegisterDbHandler(srv.Server(), new(handler.Db))

	// Run service
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
