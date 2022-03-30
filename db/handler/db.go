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
	pauth "github.com/micro/services/pkg/auth"
	gorm2 "github.com/micro/services/pkg/gorm"
	adminpb "github.com/micro/services/pkg/service/proto"
	"github.com/micro/services/pkg/tenant"
	"github.com/patrickmn/go-cache"
	"google.golang.org/protobuf/types/known/structpb"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

const idKey = "id"
const _idKey = "_id"
const stmt = "create table if not exists %v(id text not null, data jsonb, primary key(id)); alter table %v add created_at timestamptz; alter table %v add updated_at timestamptz"
const truncateStmt = `truncate table "%v"`
const dropTableStmt = `drop table "%v"`
const renameTableStmt = `ALTER TABLE "%v" RENAME TO "%v"`

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

func correctFieldName(s string, isText bool) string {
	operator := "->"
	if isText {
		// https: //stackoverflow.com/questions/27215216/postgres-how-to-convert-a-json-string-to-text
		operator = "->>"
	}
	switch s {
	// top level fields can stay top level
	case "id": // "created_at", "updated_at",  <-- these are not special fields for now
		return s
	}
	if !strings.Contains(s, ".") {
		return fmt.Sprintf("data %v '%v'", operator, s)
	}
	paths := strings.Split(s, ".")
	ret := "data"
	for i, path := range paths {
		if i == len(paths)-1 && isText {
			ret += fmt.Sprintf(" ->> '%v'", path)
			break
		}
		ret += fmt.Sprintf(" -> '%v'", path)
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
	id := req.Id

	// check the record for an id field
	if len(id) == 0 {
		// try use an id from the record
		if mid, ok := m[idKey].(string); ok {
			id = mid
		} else {
			// set id as uuid
			id = uuid.New().String()
			// inject into record
			m[idKey] = id
		}
	}

	bs, _ := json.Marshal(m)

	err = db.Table(tableName).Create(&Record{
		ID:   id,
		Data: bs,
	}).Error
	if err != nil {
		return err
	}

	// set the response id
	rsp.Id = id

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

	// if the id is blank then check the data
	if len(id) == 0 {
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
			logger.Infof("Query field: %v, op: %v, value: %v", query.Field, query.Op, query.Value)
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
			queryField := correctFieldName(query.Field, typ == "text")
			db = db.Where(fmt.Sprintf("(%v)::%v %v ?", queryField, typ, op), query.Value)
		}
	}

	orderField := "created_at"
	if req.OrderBy != "" {
		orderField = req.OrderBy
	}
	orderField = correctFieldName(orderField, false)

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
	err = db.Debug().Find(&recs).Error
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

		// only inject the ID if it does not exist
		if id, ok := ma[idKey]; !ok {
			ma[idKey] = rec.ID
		} else if id != rec.ID {
			// inject an _id key because
			// they don't match e.g user defined
			// an id field in their data
			// and separately set an id
			ma[_idKey] = rec.ID
		}

		m, _ = json.Marshal(ma)
		s := &structpb.Struct{}

		if err = s.UnmarshalJSON(m); err != nil {
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

func (e *Db) DropTable(ctx context.Context, req *db.DropTableRequest, rsp *db.DropTableResponse) error {
	tableName, err := e.tableName(ctx, req.Table)
	if err != nil {
		return err
	}
	logger.Infof("Dropping table '%v'", tableName)

	db, err := e.GetDBConn(ctx)
	if err != nil {
		return err
	}
	return db.Exec(fmt.Sprintf(dropTableStmt, tableName)).Error
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

func (e *Db) RenameTable(ctx context.Context, req *db.RenameTableRequest, rsp *db.RenameTableResponse) error {
	if req.From == "" || req.To == "" {
		return errors.BadRequest("db.renameTable", "must provide table names")
	}

	oldtableName, err := e.tableName(ctx, req.From)
	if err != nil {
		return err
	}

	newtableName, err := e.tableName(ctx, req.To)
	if err != nil {
		return err
	}

	db, err := e.GetDBConn(ctx)
	if err != nil {
		return err
	}

	stmt := fmt.Sprintf(renameTableStmt, oldtableName, newtableName)
	logger.Info(stmt)
	return db.Debug().Exec(stmt).Error
}

func (e *Db) ListTables(ctx context.Context, req *db.ListTablesRequest, rsp *db.ListTablesResponse) error {
	tenantId, ok := tenant.FromContext(ctx)
	if !ok {
		tenantId = "micro"
	}
	tenantId = strings.Replace(strings.Replace(tenantId, "/", "_", -1), "-", "_", -1)

	db, err := e.GetDBConn(ctx)
	if err != nil {
		return err
	}

	var tables []string
	if err := db.Table("information_schema.tables").Select("table_name").Where("table_schema = ?", "public").Find(&tables).Error; err != nil {
		return err
	}
	rsp.Tables = []string{}
	for _, v := range tables {
		if strings.HasPrefix(v, tenantId) {
			rsp.Tables = append(rsp.Tables, strings.Replace(v, tenantId+"_", "", -1))
		}
	}
	return nil
}

func (e *Db) DeleteData(ctx context.Context, request *adminpb.DeleteDataRequest, response *adminpb.DeleteDataResponse) error {
	method := "admin.DeleteData"
	_, err := pauth.VerifyMicroAdmin(ctx, method)
	if err != nil {
		return err
	}

	if len(request.TenantId) < 10 { // deliberate length check so we don't delete all the things
		return errors.BadRequest(method, "Missing tenant ID")
	}

	split := strings.Split(request.TenantId, "/")
	tctx := tenant.NewContext(split[1], split[0], split[1])

	tenantId := request.TenantId
	tenantId = strings.Replace(strings.Replace(tenantId, "/", "_", -1), "-", "_", -1)

	db, err := e.GetDBConn(tctx)
	if err != nil {
		return err
	}

	var tables []string
	if err := db.Table("information_schema.tables").Select("table_name").Where("table_schema = ?", "public").Find(&tables).Error; err != nil {
		return err
	}
	dropCount := 0
	for _, v := range tables {
		if !strings.HasPrefix(v, tenantId) {
			continue
		}
		if err := db.Exec(fmt.Sprintf(dropTableStmt, v)).Error; err != nil {
			return err
		}
		dropCount++
	}

	logger.Infof("Deleted %d tables for %s", dropCount, request.TenantId)
	return nil
}

func (e *Db) Usage(ctx context.Context, request *adminpb.UsageRequest, response *adminpb.UsageResponse) error {
	return nil
}
