package handler_test

import (
	"database/sql"
	"os"
	"testing"
	"time"

	"github.com/micro/services/codes/handler"
)

func testHandler(t *testing.T) *handler.Codes {
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
	if _, err := sqlDB.Exec("DROP TABLE IF EXISTS micro_codes CASCADE"); err != nil {
		t.Fatalf("Error cleaning database: %v", err)
	}

	h := &handler.Codes{Time: time.Now}
	h.DBConn(sqlDB).Migrations(&handler.Code{})
	return h
}
