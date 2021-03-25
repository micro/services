package main

import (
	"database/sql"

	"github.com/micro/services/invites/handler"
	pb "github.com/micro/services/invites/proto"

	"github.com/micro/micro/v3/service"
	"github.com/micro/micro/v3/service/config"
	"github.com/micro/micro/v3/service/logger"

	_ "github.com/jackc/pgx/v4/stdlib"
)

var dbAddress = "postgresql://postgres:postgres@localhost:5432/invites?sslmode=disable"

func main() {
	// Create service
	srv := service.New(
		service.Name("invites"),
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
	h := &handler.Invites{}
	h.DBConn(sqlDB).Migrations(&handler.Invite{})

	// Register handler
	pb.RegisterInvitesHandler(srv.Server(), h)

	// Run service
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
