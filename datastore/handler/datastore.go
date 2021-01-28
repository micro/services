package handler

import (
	"context"
	"encoding/json"
	"strings"

	log "github.com/micro/micro/v3/service/logger"
	"github.com/micro/micro/v3/service/model"

	datastore "github.com/micro/services/datastore/proto"
)

var indexIndex = model.Index{
	FieldName: "TypeOf",
}

type IndexRecord struct {
	ID     string
	TypeOf string
	Index  model.Index
}

type Datastore struct {
}

func (e *Datastore) Create(ctx context.Context, req *datastore.CreateRequest, rsp *datastore.CreateResponse) error {
	log.Info("Received Datastore.Create request")
	m := map[string]interface{}{}
	err := json.Unmarshal([]byte(req.Value), &m)
	if err != nil {
		return err
	}
	indexes, err := e.getIndexes(ctx)
	if err != nil {
		return err
	}
	db := model.New(map[string]interface{}{}, &model.Options{
		Indexes: indexes,
	})
	return db.Context(ctx).Create(m)
}

func (e *Datastore) Update(ctx context.Context, req *datastore.UpdateRequest, rsp *datastore.UpdateResponse) error {
	log.Info("Received Datastore.Update request")
	m := map[string]interface{}{}
	err := json.Unmarshal([]byte(req.Value), &m)
	if err != nil {
		return err
	}
	indexes, err := e.getIndexes(ctx)
	if err != nil {
		return err
	}
	db := model.New(map[string]interface{}{}, &model.Options{
		Indexes: indexes,
	})
	return db.Context(ctx).Update(m)
}

func (e *Datastore) getIndexes(ctx context.Context) ([]model.Index, error) {
	indexDb := model.New(map[string]interface{}{}, &model.Options{
		Indexes: []model.Index{indexIndex},
	})
	result := []IndexRecord{}
	err := indexDb.Context(ctx).Read(model.QueryEquals("TypeOf", "_index"), &result)
	if err != nil {
		return nil, err
	}
	indexes := []model.Index{}
	for _, v := range result {
		indexes = append(indexes, v.Index)
	}
	return indexes, nil
}

func (e *Datastore) Read(ctx context.Context, req *datastore.ReadRequest, rsp *datastore.ReadResponse) error {
	log.Info("Received Datastore.Read request")
	q := toQuery(req.Query)
	result := []map[string]interface{}{}
	indexes, err := e.getIndexes(ctx)
	if err != nil {
		return err
	}
	db := model.New(map[string]interface{}{}, &model.Options{
		Indexes: indexes,
	})
	err = db.Context(ctx).Read(q, &result)
	if err != nil {
		return err
	}
	js, err := json.Marshal(result)
	rsp.Value = string(js)
	return err
}

func (e *Datastore) CreateIndex(ctx context.Context, req *datastore.CreateIndexRequest, rsp *datastore.CreateIndexResponse) error {
	log.Info("Received Datastore.Index request")

	index := toIndex(req.Index)
	indexRecord := IndexRecord{
		ID:     index.FieldName + index.Type + index.Order.FieldName + string(index.Order.Type),
		Index:  index,
		TypeOf: "_index",
	}
	db := model.New(IndexRecord{}, &model.Options{
		Indexes: []model.Index{indexIndex},
	})
	return db.Context(ctx).Create(indexRecord)
}

func (e *Datastore) Delete(ctx context.Context, req *datastore.DeleteRequest, rsp *datastore.DeleteResponse) error {
	log.Info("Received Datastore.Delete request")
	q := toQuery(req.Query)
	return model.New(map[string]interface{}{}, nil).Context(ctx).Delete(q)
}

func toQuery(pquery *datastore.Query) model.Query {
	q := model.QueryEquals(pquery.Index.FieldName, pquery.Value)
	if pquery.Order != nil {
		q.Order.FieldName = pquery.Order.FieldName
		q.Order.Type = model.OrderType(pquery.Order.OrderType.String())
	}
	return q
}

func toIndex(pindex *datastore.Index) model.Index {
	i := model.Index{
		FieldName: pindex.FieldName,
		Type:      pindex.Type,
		Unique:    pindex.Unique,
	}
	if pindex.Order != nil {
		i.Order = model.Order{
			FieldName: pindex.FieldName,
			Type:      model.OrderType(strings.ToLower(pindex.Order.OrderType.String())),
		}
	}
	return i
}
