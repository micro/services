package main

import (
	"database/sql"
	"time"

	"github.com/micro/services/codes/handler"
	pb "github.com/micro/services/codes/proto"

	"github.com/micro/micro/v3/service"
	"github.com/micro/micro/v3/service/config"
	"github.com/micro/micro/v3/service/logger"

	_ "github.com/jackc/pgx/v4/stdlib"
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
	sqlDB, err := sql.Open("pgx", addr)
	if err != nil {
		logger.Fatalf("Failed to open connection to DB %s", err)
	}

	h := &handler.Codes{Time: time.Now}
	h.DBConn(sqlDB).Migrations(&handler.Code{})
	// Register handler
	pb.RegisterCodesHandler(srv.Server(), h)

	// Run service
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
