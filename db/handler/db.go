package handler

import (
	"context"
	"database/sql/driver"
	"encoding/json"

	db "github.com/micro/services/db/proto"
	gorm2 "github.com/micro/services/pkg/gorm"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type JSONB map[string]interface{}

func (j JSONB) Value() (driver.Value, error) {
	valueString, err := json.Marshal(j)
	return string(valueString), err
}

func (j *JSONB) Scan(value interface{}) error {
	if err := json.Unmarshal(value.([]byte), &j); err != nil {
		return err
	}
	return nil
}

type Record struct {
	gorm.Model
	Data JSONB `sql:"type:jsonb"`
}

type Db struct {
	gorm2.Helper
}

// Call is a single request handler called via client.Call or the generated client code
func (e *Db) Create(ctx context.Context, req *db.CreateRequest, rsp *db.CreateResponse) error {
	db, err := e.GetDBConn(ctx)
	if err != nil {
		return err
	}
	rec := &Record{}
	err = json.Unmarshal([]byte(req.Record), &rec.Data)
	if err != nil {
		return err
	}

	return db.Table(req.Table).Create(rec).Error
}

func (e *Db) Update(ctx context.Context, req *db.UpdateRequest, rsp *db.UpdateResponse) error {

	return nil
}

func (e *Db) Read(ctx context.Context, req *db.ReadRequest, rsp *db.ReadResponse) error {
	recs := []Record{}
	queries, err := Parse(req.Query)
	if err != nil {
		return err
	}
	db, err := e.GetDBConn(ctx)
	if err != nil {
		return err
	}
	gq := datatypes.JSONQuery("data")
	for _, query := range queries {
		switch query.Op {
		case itemEquals:
			gq = gq.Equals(query.Value, query.Field)
		}
	}
	err = db.Table(req.Table).Where(gq).Find(recs).Error
	if err != nil {
		return err
	}
	bts, _ := json.Marshal(recs)
	rsp.Records = string(bts)
	return nil
}

func (e *Db) Delete(ctx context.Context, req *db.DeleteRequest, rsp *db.DeleteResponse) error {

	return nil
}
