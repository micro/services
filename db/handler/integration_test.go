package handler

import (
	"context"
	"encoding/json"
	"testing"

	"database/sql"

	db "github.com/micro/services/db/proto"
	"google.golang.org/protobuf/types/known/structpb"
)

const dbAddr = "postgresql://postgres:postgres@postgres:5432/postgres?sslmode=disable"

func getHandler(t *testing.T) *Db {
	sqlDB, err := sql.Open("pgx", dbAddr)
	if err != nil {
		t.Fatalf("Failed to open connection to DB %s", err)
	}
	h := &Db{}
	h.DBConn(sqlDB)
	return h
}

func TestCreate(t *testing.T) {
	h := getHandler(t)
	record, _ := json.Marshal(map[string]interface{}{
		"name":     "Jane",
		"age":      42,
		"isActive": true,
		"id":       "1",
	})
	rec := &structpb.Struct{}
	err := rec.UnmarshalJSON(record)
	if err != nil {
		t.Fatal(err)
	}
	err = h.Create(context.TODO(), &db.CreateRequest{
		Table:  "users",
		Record: rec,
	}, &db.CreateResponse{})
	if err != nil {
		t.Fatal(err)
	}
}
