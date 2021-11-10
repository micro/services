package handler

import (
	"context"
	"encoding/json"
	"testing"

	"database/sql"

	"github.com/micro/micro/v3/service/auth"
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

func TestBasic(t *testing.T) {
	h := getHandler(t)
	ctx := auth.ContextWithAccount(context.Background(), &auth.Account{Issuer: "basic_test", ID: "test"})
	rs := []map[string]interface{}{
		{
			"name":     "Jane",
			"age":      42,
			"isActive": true,
			"id":       "1",
		},
		{
			"name":     "Joe",
			"age":      112,
			"isActive": false,
			"id":       "2",
		},
	}
	for _, v := range rs {
		record, _ := json.Marshal(v)
		rec := &structpb.Struct{}
		err := rec.UnmarshalJSON(record)
		if err != nil {
			t.Fatal(err)
		}
		err = h.Create(ctx, &db.CreateRequest{
			Table:  "users",
			Record: rec,
		}, &db.CreateResponse{})
		if err != nil {
			t.Fatal(err)
		}
	}

	t.Run("number ==", func(t *testing.T) {
		readRsp := &db.ReadResponse{}
		err := h.Read(ctx, &db.ReadRequest{
			Table: "users",
			Query: "age == 112",
		}, readRsp)
		if err != nil {
			t.Fatal(err)
		}
		if len(readRsp.Records) != 1 || readRsp.Records[0].AsMap()["id"].(string) != "2" {
			t.Fatal(readRsp)
		}
	})

	t.Run("number <", func(t *testing.T) {
		readRsp := &db.ReadResponse{}
		err := h.Read(ctx, &db.ReadRequest{
			Table: "users",
			Query: "age < 100",
		}, readRsp)
		if err != nil {
			t.Fatal(err)
		}
		if len(readRsp.Records) != 1 || readRsp.Records[0].AsMap()["id"].(string) != "1" {
			t.Fatal(readRsp)
		}
	})

	t.Run("number >", func(t *testing.T) {
		readRsp := &db.ReadResponse{}
		err := h.Read(ctx, &db.ReadRequest{
			Table: "users",
			Query: "age > 100",
		}, readRsp)
		if err != nil {
			t.Fatal(err)
		}
		if len(readRsp.Records) != 1 || readRsp.Records[0].AsMap()["id"].(string) != "2" {
			t.Fatal(readRsp)
		}
	})

	t.Run("number !=", func(t *testing.T) {
		readRsp := &db.ReadResponse{}
		err := h.Read(ctx, &db.ReadRequest{
			Table: "users",
			Query: "age != 42",
		}, readRsp)
		if err != nil {
			t.Fatal(err)
		}
		if len(readRsp.Records) != 1 || readRsp.Records[0].AsMap()["id"].(string) != "2" {
			t.Fatal(readRsp)
		}
	})

	t.Run("bool ==", func(t *testing.T) {
		readRsp := &db.ReadResponse{}
		err := h.Read(ctx, &db.ReadRequest{
			Table: "users",
			Query: "isActive == false",
		}, readRsp)
		if err != nil {
			t.Fatal(err)
		}
		if len(readRsp.Records) != 1 || readRsp.Records[0].AsMap()["id"].(string) != "2" {
			t.Fatal(readRsp)
		}
	})

	t.Run("bool !=", func(t *testing.T) {
		readRsp := &db.ReadResponse{}
		err := h.Read(ctx, &db.ReadRequest{
			Table: "users",
			Query: "isActive != false",
		}, readRsp)
		if err != nil {
			t.Fatal(err)
		}
		if len(readRsp.Records) != 1 || readRsp.Records[0].AsMap()["id"].(string) != "1" {
			t.Fatal(readRsp)
		}
	})

	t.Run("string ==", func(t *testing.T) {
		readRsp := &db.ReadResponse{}
		err := h.Read(ctx, &db.ReadRequest{
			Table: "users",
			Query: "name == 'Jane'",
		}, readRsp)
		if err != nil {
			t.Fatal(err)
		}
		if len(readRsp.Records) != 1 || readRsp.Records[0].AsMap()["id"].(string) != "1" {
			t.Fatal(readRsp)
		}
	})

	t.Run("string !=", func(t *testing.T) {
		readRsp := &db.ReadResponse{}
		err := h.Read(ctx, &db.ReadRequest{
			Table: "users",
			Query: "name != 'Jane'",
		}, readRsp)
		if err != nil {
			t.Fatal(err)
		}
		if len(readRsp.Records) != 1 || readRsp.Records[0].AsMap()["id"].(string) != "2" {
			t.Fatal(readRsp)
		}
	})

	t.Run("order number asc", func(t *testing.T) {
		readRsp := &db.ReadResponse{}
		err := h.Read(ctx, &db.ReadRequest{
			Table:   "users",
			OrderBy: "age",
			Order:   "asc",
		}, readRsp)
		if err != nil {
			t.Fatal(err)
		}
		if len(readRsp.Records) != 2 || readRsp.Records[0].AsMap()["id"].(string) != "1" || readRsp.Records[1].AsMap()["id"].(string) != "2" {
			t.Fatal(readRsp)
		}
	})

	t.Run("order number desc", func(t *testing.T) {
		readRsp := &db.ReadResponse{}
		err := h.Read(ctx, &db.ReadRequest{
			Table:   "users",
			OrderBy: "age",
			Order:   "desc",
		}, readRsp)
		if err != nil {
			t.Fatal(err)
		}
		if len(readRsp.Records) != 2 || readRsp.Records[0].AsMap()["id"].(string) != "2" || readRsp.Records[1].AsMap()["id"].(string) != "1" {
			t.Fatal(readRsp)
		}
	})

	t.Run("order number desc, limit", func(t *testing.T) {
		readRsp := &db.ReadResponse{}
		err := h.Read(ctx, &db.ReadRequest{
			Table:   "users",
			OrderBy: "age",
			Order:   "desc",
			Limit:   1,
		}, readRsp)
		if err != nil {
			t.Fatal(err)
		}
		if len(readRsp.Records) != 1 || readRsp.Records[0].AsMap()["id"].(string) != "2" {
			t.Fatal(readRsp)
		}
	})

	t.Run("order number desc, limit, offset", func(t *testing.T) {
		readRsp := &db.ReadResponse{}
		err := h.Read(ctx, &db.ReadRequest{
			Table:   "users",
			OrderBy: "age",
			Order:   "desc",
			Limit:   1,
			Offset:  1,
		}, readRsp)
		if err != nil {
			t.Fatal(err)
		}
		if len(readRsp.Records) != 1 || readRsp.Records[0].AsMap()["id"].(string) != "1" {
			t.Fatal(readRsp)
		}
	})
}
