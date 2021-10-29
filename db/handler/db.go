package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/micro/micro/v3/service/errors"
	"github.com/micro/micro/v3/service/logger"
	db "github.com/micro/services/db/proto"
	gorm2 "github.com/micro/services/pkg/gorm"
	"github.com/micro/services/pkg/tenant"
	"github.com/patrickmn/go-cache"
	"google.golang.org/protobuf/types/known/structpb"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

const idKey = "id"
const stmt = "create table if not exists %v(id text not null, data jsonb, primary key(id)); alter table %v add created_at timestamptz; alter table %v add updated_at timestamptz"
const truncateStmt = `truncate table "%v"`

var re = regexp.MustCompile("^[a-zA-Z0-9_]*$")
var c = cache.New(5*time.Minute, 10*time.Minute)

type Record struct {
	ID   string
	Data datatypes.JSON `json:"data"`
	// private field, ignored from gorm
	table     string `gorm:"-"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Db struct {
	gorm2.Helper
}

func correctFieldName(s string) string {
	switch s {
	// top level fields can stay top level
	case "id": // "created_at", "updated_at",  <-- these are not special fields for now
		return s
	}
	if !strings.Contains(s, ".") {
		return fmt.Sprintf("data ->> '%v'", s)
	}
	paths := strings.Split(s, ".")
	ret := "data"
	for _, path := range paths {
		ret += fmt.Sprintf(" ->> '%v'", path)
	}
	return ret
}

func (e *Db) tableName(ctx context.Context, t string) (string, error) {
	tenantId, ok := tenant.FromContext(ctx)
	if !ok {
		tenantId = "micro"
	}
	if t == "" {
		t = "default"
	}
	t = strings.ToLower(t)
	t = strings.Replace(t, "-", "_", -1)
	tenantId = strings.Replace(strings.Replace(tenantId, "/", "_", -1), "-", "_", -1)

	tableName := tenantId + "_" + t
	if !re.Match([]byte(tableName)) {
		return "", fmt.Errorf("table name %v is invalid", t)
	}

	return tableName, nil
}

// Call is a single request handler called via client.Call or the generated client code
func (e *Db) Create(ctx context.Context, req *db.CreateRequest, rsp *db.CreateResponse) error {
	if len(req.Record.AsMap()) == 0 {
		return errors.BadRequest("db.create", "missing record")
	}

	tableName, err := e.tableName(ctx, req.Table)
	if err != nil {
		return err
	}
	logger.Infof("Inserting into table '%v'", tableName)

	db, err := e.GetDBConn(ctx)
	if err != nil {
		return err
	}
	_, ok := c.Get(tableName)
	if !ok {
		logger.Infof("Creating table '%v'", tableName)
		db.Exec(fmt.Sprintf(stmt, tableName, tableName, tableName))
		c.Set(tableName, true, 0)
	}

	m := req.Record.AsMap()
	if _, ok := m[idKey].(string); !ok {
		m[idKey] = uuid.New().String()
	}
	bs, _ := json.Marshal(m)

	err = db.Table(tableName).Create(&Record{
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
	if len(req.Record.AsMap()) == 0 {
		return errors.BadRequest("db.update", "missing record")
	}
	tableName, err := e.tableName(ctx, req.Table)
	if err != nil {
		return err
	}
	logger.Infof("Updating table '%v'", tableName)

	db, err := e.GetDBConn(ctx)
	if err != nil {
		return err
	}

	m := req.Record.AsMap()

	// where ID is specified do a single update record update
	id := req.Id
	if v, ok := m[idKey].(string); ok && id == "" {
		id = v
	}

	// if the id is blank then check the data
	if len(req.Id) == 0 {
		var ok bool
		id, ok = m[idKey].(string)
		if !ok {
			return fmt.Errorf("update failed: missing id")
		}
	}

	return db.Transaction(func(tx *gorm.DB) error {
		rec := []Record{}
		err = tx.Table(tableName).Where("id = ?", id).Find(&rec).Error
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
		for k, v := range m {
			old[k] = v
		}
		bs, _ := json.Marshal(old)

		return tx.Table(tableName).Save(&Record{
			ID:   id,
			Data: bs,
		}).Error
	})
}

func (e *Db) Read(ctx context.Context, req *db.ReadRequest, rsp *db.ReadResponse) error {
	recs := []Record{}
	queries, err := Parse(req.Query)
	if err != nil {
		return err
	}
	tableName, err := e.tableName(ctx, req.Table)
	if err != nil {
		return err
	}

	db, err := e.GetDBConn(ctx)
	if err != nil {
		return err
	}
	_, ok := c.Get(tableName)
	if !ok {
		logger.Infof("Creating table '%v'", tableName)
		db.Exec(fmt.Sprintf(stmt, tableName, tableName, tableName))
		c.Set(tableName, true, 0)
	}

	if req.Limit > 1000 {
		return errors.BadRequest("db.read", fmt.Sprintf("limit over 1000 is invalid, you specified %v", req.Limit))
	}
	if req.Limit == 0 {
		req.Limit = 25
	}

	db = db.Table(tableName)
	if req.Id != "" {
		logger.Infof("Query by id: %v", req.Id)
		db = db.Where("id = ?", req.Id)
	} else {
		for _, query := range queries {
			logger.Infof("Query field: %v, op: %v, type: %v", query.Field, query.Op, query.Value)
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
			queryField := correctFieldName(query.Field)
			db = db.Where(fmt.Sprintf("(%v)::%v %v ?", queryField, typ, op), query.Value)
		}
	}

	orderField := "created_at"
	if req.OrderBy != "" {
		orderField = req.OrderBy
	}
	orderField = correctFieldName(orderField)

	ordering := "asc"
	if req.Order != "" {
		switch strings.ToLower(req.Order) {
		case "asc":
			ordering = "asc"
		case "", "desc":
			ordering = "desc"
		default:
			return errors.BadRequest("db.read", "invalid ordering: "+req.Order)
		}
	}

	db = db.Order(orderField + " " + ordering).Offset(int(req.Offset)).Limit(int(req.Limit))
	err = db.Find(&recs).Error
	if err != nil {
		return err
	}

	rsp.Records = []*structpb.Struct{}
	for _, rec := range recs {
		m, err := rec.Data.MarshalJSON()
		if err != nil {
			return err
		}
		ma := map[string]interface{}{}
		json.Unmarshal(m, &ma)
		ma[idKey] = rec.ID
		m, _ = json.Marshal(ma)
		s := &structpb.Struct{}
		err = s.UnmarshalJSON(m)
		if err != nil {
			return err
		}
		rsp.Records = append(rsp.Records, s)
	}

	return nil
}

func (e *Db) Delete(ctx context.Context, req *db.DeleteRequest, rsp *db.DeleteResponse) error {
	if len(req.Id) == 0 {
		return errors.BadRequest("db.delete", "missing id")
	}
	tableName, err := e.tableName(ctx, req.Table)
	if err != nil {
		return err
	}
	logger.Infof("Deleting from table '%v'", tableName)

	db, err := e.GetDBConn(ctx)
	if err != nil {
		return err
	}

	return db.Table(tableName).Delete(Record{
		ID: req.Id,
	}).Error
}

func (e *Db) Truncate(ctx context.Context, req *db.TruncateRequest, rsp *db.TruncateResponse) error {
	tableName, err := e.tableName(ctx, req.Table)
	if err != nil {
		return err
	}
	logger.Infof("Truncating table '%v'", tableName)

	db, err := e.GetDBConn(ctx)
	if err != nil {
		return err
	}
	return db.Exec(fmt.Sprintf(truncateStmt, tableName)).Error
}

func (e *Db) Count(ctx context.Context, req *db.CountRequest, rsp *db.CountResponse) error {
	if req.Table == "" {
		req.Table = "default"
	}

	tableName, err := e.tableName(ctx, req.Table)
	if err != nil {
		return err
	}

	db, err := e.GetDBConn(ctx)
	if err != nil {
		return err
	}

	var a int64
	err = db.Table(tableName).Model(Record{}).Count(&a).Error
	if err != nil {
		return err
	}
	rsp.Count = int32(a)
	return nil
}
