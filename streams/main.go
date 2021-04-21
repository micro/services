package main

import (
	"database/sql"
	"time"

	"github.com/micro/services/streams/handler"
	pb "github.com/micro/services/streams/proto"

	"github.com/micro/micro/v3/service"
	"github.com/micro/micro/v3/service/config"
	"github.com/micro/micro/v3/service/events"
	"github.com/micro/micro/v3/service/logger"
)

var dbAddress = "postgresql://postgres:postgres@localhost:5432/streams?sslmode=disable"

func main() {
	// Create service
	srv := service.New(
		service.Name("streams"),
		service.Version("latest"),
	)

	// Connect to the database
	cfg, err := config.Get("streams.database")
	if err != nil {
		logger.Fatalf("Error loading config: %v", err)
	}
	addr := cfg.String(dbAddress)
	sqlDB, err := sql.Open("pgx", addr)
	if err != nil {
		logger.Fatalf("Failed to open connection to DB %s", err)
	}

	h := &handler.Streams{
		Events: events.DefaultStream,
		Time:   time.Now,
	}
	h.DBConn(sqlDB).Migrations(&handler.Token{})

	// Register handler
	pb.RegisterStreamsHandler(srv.Server(), h)

	// Run service
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
