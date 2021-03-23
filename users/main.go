package main

import (
	"time"

	"github.com/micro/micro/v3/service"
	"github.com/micro/micro/v3/service/config"
	"github.com/micro/micro/v3/service/logger"
	"github.com/micro/services/users/handler"
	pb "github.com/micro/services/users/proto"
	"gorm.io/driver/postgres"
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
	// Register handler
	pb.RegisterUsersHandler(srv.Server(), handler.NewHandler(time.Now, postgres.Open(addr)))

	// Run service
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
