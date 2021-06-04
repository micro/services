package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/micro/micro/v3/service/errors"
	"github.com/micro/micro/v3/service/logger"
	db "github.com/micro/services/db/proto"
	gorm2 "github.com/micro/services/pkg/gorm"
	"github.com/micro/services/pkg/tenant"
	"github.com/patrickmn/go-cache"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

const idKey = "id"
const stmt = "create table %v(id text not null, data jsonb, primary key(id));"

var c = cache.New(5*time.Minute, 10*time.Minute)

type Record struct {
	ID   string
	Data datatypes.JSON `json:"data"`
	// private field, ignored from gorm
	table string `gorm:"-"`
}

type Db struct {
	gorm2.Helper
}

// Call is a single request handler called via client.Call or the generated client code
func (e *Db) Create(ctx context.Context, req *db.CreateRequest, rsp *db.CreateResponse) error {
	if len(req.Record) == 0 {
		return errors.BadRequest("db.create", "missing record")
	}
	tenantId, ok := tenant.FromContext(ctx)
	if !ok {
		tenantId = "micro"
	}
	tenantId = strings.Replace(tenantId, "/", "_", -1)
	db, err := e.GetDBConn(ctx)
	if err != nil {
		return err
	}
	_, ok = c.Get(req.Table)
	if !ok {
		db.Exec(fmt.Sprintf(stmt, tenantId+"_"+req.Table))
		logger.Info(tenantId + "_" + req.Table)
		c.Set(req.Table, true, 0)
	}

	m := map[string]interface{}{}
	err = json.Unmarshal([]byte(req.Record), &m)
	if err != nil {
		return err
	}
	if _, ok := m[idKey].(string); !ok {
		m[idKey] = uuid.New().String()
	}
	bs, _ := json.Marshal(m)

	err = db.Table(tenantId + "_" + req.Table).Create(Record{
		ID:   m[idKey].(string),
		Data: bs,
	}).Error
	if err != nil {
		return err
	}

	// set the response id
	rsp.Id = m[idKey].(string)

	return nil
}

func (e *Db) Update(ctx context.Context, req *db.UpdateRequest, rsp *db.UpdateResponse) error {
	if len(req.Record) == 0 {
		return errors.BadRequest("db.update", "missing record")
	}
	tenantId, ok := tenant.FromContext(ctx)
	if !ok {
		tenantId = "micro"
	}
	tenantId = strings.Replace(tenantId, "/", "_", -1)
	db, err := e.GetDBConn(ctx)
	if err != nil {
		return err
	}

	m := map[string]interface{}{}
	err = json.Unmarshal([]byte(req.Record), &m)
	if err != nil {
		return err
	}

	// where ID is specified do a single update record update
	id, ok := m[idKey].(string)
	if !ok {
		return fmt.Errorf("update failed: missing id")
	}

	db.Transaction(func(tx *gorm.DB) error {
		rec := []Record{}
		err = tx.Table(tenantId+"_"+req.Table).Where("ID = ?", id).Find(&rec).Error
		if err != nil {
			return err
		}
		if len(rec) == 0 {
			return fmt.Errorf("update failed: not found")
		}
		old := map[string]interface{}{}
		err = json.Unmarshal(rec[0].Data, &old)
		if err != nil {
			return err
		}
		for k, v := range old {
			m[k] = v
		}
		bs, _ := json.Marshal(m)

		return tx.Table(tenantId + "_" + req.Table).Save(Record{
			ID:   m[idKey].(string),
			Data: bs,
		}).Error
	})
	return nil
}

func (e *Db) Read(ctx context.Context, req *db.ReadRequest, rsp *db.ReadResponse) error {
	recs := []Record{}
	queries, err := Parse(req.Query)
	if err != nil {
		return err
	}
	tenantId, ok := tenant.FromContext(ctx)
	if !ok {
		tenantId = "micro"
	}
	tenantId = strings.Replace(tenantId, "/", "_", -1)
	db, err := e.GetDBConn(ctx)
	if err != nil {
		return err
	}
	db = db.Table(tenantId + "_" + req.Table)
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
		ma[idKey] = rec.ID
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
