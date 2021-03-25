package main

import (
	"database/sql"
	"time"

	"github.com/micro/services/threads/handler"
	pb "github.com/micro/services/threads/proto"

	"github.com/micro/micro/v3/service"
	"github.com/micro/micro/v3/service/config"
	"github.com/micro/micro/v3/service/logger"

	_ "github.com/jackc/pgx/v4/stdlib"
)

var dbAddress = "postgresql://postgres:postgres@localhost:5432/threads?sslmode=disable"

func main() {
	// Create service
	srv := service.New(
		service.Name("threads"),
		service.Version("latest"),
	)

	// Connect to the database
	cfg, err := config.Get("threads.database")
	if err != nil {
		logger.Fatalf("Error loading config: %v", err)
	}
	addr := cfg.String(dbAddress)
	sqlDB, err := sql.Open("pgx", addr)
	if err != nil {
		logger.Fatalf("Failed to open connection to DB %s", err)
	}

	h := &handler.Threads{Time: time.Now}
	h.DBConn(sqlDB).Migrations(&handler.Conversation{}, &handler.Message{})
	// Register handler
	pb.RegisterThreadsHandler(srv.Server(), h)

	// Run service
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
