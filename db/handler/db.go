package handler

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/micro/micro/v3/service/errors"
	db "github.com/micro/services/db/proto"
	gorm2 "github.com/micro/services/pkg/gorm"
	"gorm.io/datatypes"
)

type Record struct {
	ID   string
	Data datatypes.JSON `json:"data"`
}

type Db struct {
	gorm2.Helper
}

// Call is a single request handler called via client.Call or the generated client code
func (e *Db) Create(ctx context.Context, req *db.CreateRequest, rsp *db.CreateResponse) error {
	if len(req.Record) == 0 {
		return errors.BadRequest("db.create", "missing record")
	}

	db, err := e.GetDBConn(ctx)
	if err != nil {
		return err
	}
	m := map[string]interface{}{}
	err = json.Unmarshal([]byte(req.Record), &m)
	if err != nil {
		return err
	}
	if _, ok := m["ID"].(string); !ok {
		m["ID"] = uuid.New().String()
	}
	bs, _ := json.Marshal(m)

	err = db.Table(req.Table).Create(Record{
		ID:   m["ID"].(string),
		Data: bs,
	}).Error
	if err != nil {
		return err
	}

	// set the response id
	rsp.Id = m["ID"].(string)

	return nil
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
	db = db.Table(req.Table)
	for _, query := range queries {
		typ := "text"
		switch query.Value.(type) {
		case int64:
			typ = "int"
		case bool:
			typ = "boolean"
		}
		op := ""
		switch query.Op {
		case itemEquals:
			op = "="
		case itemGreaterThan:
			op = ">"
		case itemGreaterThanEquals:
			op = ">="
		case itemLessThan:
			op = "<"
		case itemLessThanEquals:
			op = "<="
		case itemNotEquals:
			op = "!="
		}
		db = db.Where(fmt.Sprintf("(data ->> '%v')::%v %v ?", query.Field, typ, op), query.Value)
	}
	err = db.Find(&recs).Error
	if err != nil {
		return err
	}
	ret := []map[string]interface{}{}
	for _, rec := range recs {
		m, err := rec.Data.MarshalJSON()
		if err != nil {
			return err
		}
		ma := map[string]interface{}{}
		json.Unmarshal(m, &ma)
		ma["ID"] = rec.ID
		ret = append(ret, ma)
	}
	bs, _ := json.Marshal(ret)
	rsp.Records = string(bs)
	return nil
}

func (e *Db) Delete(ctx context.Context, req *db.DeleteRequest, rsp *db.DeleteResponse) error {
	if len(req.Id) == 0 {
		return errors.BadRequest("db.delete", "missing id")
	}

	db, err := e.GetDBConn(ctx)
	if err != nil {
		return err
	}

	return db.Table(req.Table).Delete(Record{
		ID: req.Id,
	}).Error
}
