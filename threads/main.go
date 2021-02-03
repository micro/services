package main

import (
	"time"

	"github.com/micro/services/threads/handler"
	pb "github.com/micro/services/threads/proto"

	"github.com/micro/micro/v3/service"
	"github.com/micro/micro/v3/service/config"
	"github.com/micro/micro/v3/service/logger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var dbAddress = "postgresql://postgres@localhost:5432/threads?sslmode=disable"

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
	db, err := gorm.Open(postgres.Open(addr), &gorm.Config{})
	if err != nil {
		logger.Fatalf("Error connecting to database: %v", err)
	}
	if err := db.AutoMigrate(&handler.Conversation{}, &handler.Message{}); err != nil {
		logger.Fatalf("Error migrating database: %v", err)
	}

	// Register handler
	pb.RegisterThreadsHandler(srv.Server(), &handler.Threads{DB: db, Time: time.Now})

	// Run service
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
