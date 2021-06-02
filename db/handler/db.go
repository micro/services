package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/micro/micro/v3/service/logger"
	db "github.com/micro/services/db/proto"
	gorm2 "github.com/micro/services/pkg/gorm"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Record struct {
	gorm.Model
	Data datatypes.JSON `json:"data"`
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
	logger.Info(rec.Data)
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
	ret := []string{}
	for _, rec := range recs {
		m, err := rec.Data.MarshalJSON()
		if err != nil {
			return err
		}
		ret = append(ret, string(m))
	}
	rsp.Records = "[" + strings.Join(ret, ",") + "]"
	return nil
}

func (e *Db) Delete(ctx context.Context, req *db.DeleteRequest, rsp *db.DeleteResponse) error {

	return nil
}
