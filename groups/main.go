package main

import (
	"database/sql"

	"github.com/micro/micro/v3/service"
	"github.com/micro/micro/v3/service/config"
	"github.com/micro/micro/v3/service/logger"
	"github.com/micro/services/groups/handler"
	pb "github.com/micro/services/groups/proto"
)

var dbAddress = "postgresql://postgres:postgres@localhost:5432/groups?sslmode=disable"

func main() {
	// Create service
	srv := service.New(
		service.Name("groups"),
		service.Version("latest"),
	)

	// Connect to the database
	cfg, err := config.Get("groups.database")
	if err != nil {
		logger.Fatalf("Error loading config: %v", err)
	}
	addr := cfg.String(dbAddress)
	sqlDB, err := sql.Open("pgx", addr)
	if err != nil {
		logger.Fatalf("Failed to open connection to DB %s", err)
	}
	h := &handler.Groups{}
	h.Migrations(&handler.Group{}, &handler.Membership{}).DBConn(sqlDB)
	// Register handler
	pb.RegisterGroupsHandler(srv.Server(), h)

	// Run service
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
