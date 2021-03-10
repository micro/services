package main

import (
	"time"

	"github.com/micro/services/codes/handler"
	pb "github.com/micro/services/codes/proto"

	"github.com/micro/micro/v3/service"
	"github.com/micro/micro/v3/service/config"
	"github.com/micro/micro/v3/service/logger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var dbAddress = "postgresql://postgres:postgres@localhost:5432/codes?sslmode=disable"

func main() {
	// Create service
	srv := service.New(
		service.Name("codes"),
		service.Version("latest"),
	)

	// Connect to the database
	cfg, err := config.Get("codes.database")
	if err != nil {
		logger.Fatalf("Error loading config: %v", err)
	}
	addr := cfg.String(dbAddress)
	db, err := gorm.Open(postgres.Open(addr), &gorm.Config{})
	if err != nil {
		logger.Fatalf("Error connecting to database: %v", err)
	}
	if err := db.AutoMigrate(&handler.Code{}); err != nil {
		logger.Fatalf("Error migrating database: %v", err)
	}

	// Register handler
	pb.RegisterCodesHandler(srv.Server(), &handler.Codes{DB: db, Time: time.Now})

	// Run service
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
