package handler_test

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/micro/services/users/handler"
	pb "github.com/micro/services/users/proto"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func testHandler(t *testing.T) *handler.Users {
	// connect to the database
	addr := os.Getenv("POSTGRES_URL")
	if len(addr) == 0 {
		addr = "postgresql://postgres@localhost:5432/postgres?sslmode=disable"
	}
	db, err := gorm.Open(postgres.Open(addr), &gorm.Config{})
	if err != nil {
		t.Fatalf("Error connecting to database: %v", err)
	}

	// clean any data from a previous run
	if err := db.Exec("DROP TABLE IF EXISTS users, tokens CASCADE").Error; err != nil {
		t.Fatalf("Error cleaning database: %v", err)
	}

	// migrate the database
	if err := db.AutoMigrate(&handler.User{}, &handler.Token{}); err != nil {
		t.Fatalf("Error migrating database: %v", err)
	}

	return &handler.Users{DB: db, Time: time.Now}
}

func assertUsersMatch(t *testing.T, exp, act *pb.User) {
	if act == nil {
		t.Error("No user returned")
		return
	}
	assert.Equal(t, exp.Id, act.Id)
	assert.Equal(t, exp.FirstName, act.FirstName)
	assert.Equal(t, exp.LastName, act.LastName)
	assert.Equal(t, exp.Email, act.Email)
}
