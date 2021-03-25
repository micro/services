package handler_test

import (
	"context"
	"database/sql"
	"os"
	"testing"
	"time"

	"github.com/micro/micro/v3/service/auth"
	"github.com/stretchr/testify/assert"

	"github.com/micro/services/users/handler"
	pb "github.com/micro/services/users/proto"

	_ "github.com/jackc/pgx/v4/stdlib"
)

func testHandler(t *testing.T) *handler.Users {
	// connect to the database
	addr := os.Getenv("POSTGRES_URL")
	if len(addr) == 0 {
		addr = "postgresql://postgres@localhost:5432/postgres?sslmode=disable"
	}
	sqlDB, err := sql.Open("pgx", addr)
	if err != nil {
		t.Fatalf("Failed to open connection to DB %s", err)
	}
	// clean any data from a previous run
	if _, err := sqlDB.Exec("DROP TABLE IF EXISTS micro_users, micro_tokens CASCADE"); err != nil {
		t.Fatalf("Error cleaning database: %v", err)
	}

	h := handler.NewHandler(time.Now)
	h.DBConn(sqlDB).Migrations(&handler.User{}, &handler.Token{})
	return h
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
