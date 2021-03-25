package main

import (
	"database/sql"
	"time"

	"github.com/micro/services/chats/handler"
	pb "github.com/micro/services/chats/proto"

	"github.com/micro/micro/v3/service"
	"github.com/micro/micro/v3/service/config"
	"github.com/micro/micro/v3/service/logger"

	_ "github.com/jackc/pgx/v4/stdlib"
)

var dbAddress = "postgresql://postgres:postgres@localhost:5432/chats?sslmode=disable"

func main() {
	// Create service
	srv := service.New(
		service.Name("chats"),
		service.Version("latest"),
	)

	// Connect to the database
	cfg, err := config.Get("chats.database")
	if err != nil {
		logger.Fatalf("Error loading config: %v", err)
	}
	addr := cfg.String(dbAddress)
	sqlDB, err := sql.Open("pgx", addr)
	if err != nil {
		logger.Fatalf("Failed to open connection to DB %s", err)
	}

	h := &handler.Chats{Time: time.Now}
	h.DBConn(sqlDB).Migrations(&handler.Chat{}, &handler.Message{})
	// Register handler
	pb.RegisterChatsHandler(srv.Server(), h)

	// Run service
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
