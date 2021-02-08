package main

import (
	"time"

	"github.com/micro/services/streams/handler"
	pb "github.com/micro/services/streams/proto"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/micro/micro/v3/service"
	"github.com/micro/micro/v3/service/config"
	"github.com/micro/micro/v3/service/events"
	"github.com/micro/micro/v3/service/logger"
)

var dbAddress = "postgresql://postgres@localhost:5432/streams?sslmode=disable"

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
	db, err := gorm.Open(postgres.Open(addr), &gorm.Config{})
	if err != nil {
		logger.Fatalf("Error connecting to database: %v", err)
	}
	if err := db.AutoMigrate(&handler.Token{}); err != nil {
		logger.Fatalf("Error migrating database: %v", err)
	}

	// Register handler
	pb.RegisterStreamsHandler(srv.Server(), &handler.Streams{
		DB:     db,
		Events: events.DefaultStream,
		Time:   time.Now,
	})

	// Run service
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
