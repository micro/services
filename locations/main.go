package main

import (
	geo "github.com/hailocab/go-geoindex"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/micro/services/locations/handler"
	"github.com/micro/services/locations/model"
	pb "github.com/micro/services/locations/proto"

	"github.com/micro/micro/v3/service"
	"github.com/micro/micro/v3/service/config"
	"github.com/micro/micro/v3/service/logger"
)

var dbAddress = "postgresql://postgres@localhost:5432/locations?sslmode=disable"

func main() {
	// Create service
	srv := service.New(
		service.Name("locations"),
		service.Version("latest"),
	)

	// Connect to the database
	cfg, err := config.Get("locations.database")
	if err != nil {
		logger.Fatalf("Error loading config: %v", err)
	}
	addr := cfg.String(dbAddress)
	db, err := gorm.Open(postgres.Open(addr), &gorm.Config{})
	if err != nil {
		logger.Fatalf("Error connecting to database: %v", err)
	}

	// Migrate the database
	if err := db.AutoMigrate(&model.Location{}); err != nil {
		logger.Fatalf("Error migrating database: %v", err)
	}

	// Register handler
	pb.RegisterLocationsHandler(srv.Server(), &handler.Locations{
		DB: db, Geoindex: geo.NewPointsIndex(geo.Km(0.1)),
	})

	// Run service
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
