package handler_test

import (
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
	db, err := gorm.Open(postgres.Open("postgresql://postgres@localhost:5432/users?sslmode=disable"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Error connecting to database: %v", err)
	}

	// migrate the database
	if err := db.AutoMigrate(&handler.User{}, &handler.Token{}); err != nil {
		t.Fatalf("Error migrating database: %v", err)
	}

	// clean any data from a previous run
	if err := db.Exec("TRUNCATE TABLE users, tokens CASCADE").Error; err != nil {
		t.Fatalf("Error cleaning database: %v", err)
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
