package handler_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/micro/micro/v3/service/auth"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm/schema"

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
	dial := postgres.Open(addr)
	db, err := gorm.Open(dial, &gorm.Config{
		NamingStrategy: schema.NamingStrategy{TablePrefix: "micro_"},
	})
	if err != nil {
		t.Fatalf("Error connecting to database: %v", err)
	}

	// clean any data from a previous run
	if err := db.Exec("DROP TABLE IF EXISTS micro_users, micro_tokens CASCADE").Error; err != nil {
		t.Fatalf("Error cleaning database: %v", err)
	}

	// migrate the database
	if err := db.AutoMigrate(&handler.User{}, &handler.Token{}); err != nil {
		t.Fatalf("Error migrating database: %v", err)
	}
	return handler.NewHandler(time.Now, dial)
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

func microAccountCtx() context.Context {
	return auth.ContextWithAccount(context.TODO(), &auth.Account{
		Issuer: "micro",
	})
}
