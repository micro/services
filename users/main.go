package main

import (
	"time"

	"github.com/micro/services/users/handler"
	pb "github.com/micro/services/users/proto"

	"github.com/micro/micro/v3/service"
	"github.com/micro/micro/v3/service/config"
	"github.com/micro/micro/v3/service/logger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var dbAddress = "postgresql://postgres:postgres@localhost:5432/users?sslmode=disable"

func main() {
	// Create service
	srv := service.New(
		service.Name("users"),
		service.Version("latest"),
	)

	// Connect to the database
	cfg, err := config.Get("users.database")
	if err != nil {
		logger.Fatalf("Error loading config: %v", err)
	}
	addr := cfg.String(dbAddress)
	db, err := gorm.Open(postgres.Open(addr), &gorm.Config{})
	if err != nil {
		logger.Fatalf("Error connecting to database: %v", err)
	}
	if err := db.AutoMigrate(&handler.User{}, &handler.Token{}); err != nil {
		logger.Fatalf("Error migrating database: %v", err)
	}

	// Register handler
	pb.RegisterUsersHandler(srv.Server(), &handler.Users{DB: db, Time: time.Now})

	// Run service
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
