package main

import (
	"github.com/micro/services/seen/handler"
	pb "github.com/micro/services/seen/proto"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/micro/micro/v3/service"
	"github.com/micro/micro/v3/service/config"
	"github.com/micro/micro/v3/service/logger"
)

var dbAddress = "postgresql://postgres@localhost:5432/seen?sslmode=disable"

func main() {
	// Create service
	srv := service.New(
		service.Name("seen"),
		service.Version("latest"),
	)

	// Connect to the database
	cfg, err := config.Get("seen.database")
	if err != nil {
		logger.Fatalf("Error loading config: %v", err)
	}
	addr := cfg.String(dbAddress)
	db, err := gorm.Open(postgres.Open(addr), &gorm.Config{})
	if err != nil {
		logger.Fatalf("Error connecting to database: %v", err)
	}
	if err := db.AutoMigrate(&handler.SeenInstance{}); err != nil {
		logger.Fatalf("Error migrating database: %v", err)
	}

	// Register handler
	pb.RegisterSeenHandler(srv.Server(), &handler.Seen{DB: db})

	// Run service
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
