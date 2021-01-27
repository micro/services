package handler

import (
	"context"
	"encoding/json"

	log "github.com/micro/micro/v3/service/logger"
	"github.com/micro/micro/v3/service/model"

	datastore "github.com/micro/services/datastore/proto"
)

type Datastore struct {
}

func (e *Datastore) Create(ctx context.Context, req *datastore.CreateRequest, rsp *datastore.CreateResponse) error {
	log.Info("Received Datastore.Create request")
	m := map[string]interface{}{}
	err := json.Unmarshal([]byte(req.Value), &m)
	if err != nil {
		return err
	}

	return model.New(map[string]interface{}{}, nil).Context(ctx).Create(m)
}

func (e *Datastore) Update(ctx context.Context, req *datastore.UpdateRequest, rsp *datastore.UpdateResponse) error {
	log.Info("Received Datastore.Update request")
	m := map[string]interface{}{}
	err := json.Unmarshal([]byte(req.Value), &m)
	if err != nil {
		return err
	}

	return model.New(map[string]interface{}{}, nil).Context(ctx).Update(m)
}

func (e *Datastore) Read(ctx context.Context, req *datastore.ReadRequest, rsp *datastore.ReadResponse) error {
	log.Info("Received Datastore.Read request")
	q := toQuery(req.Query)
	result := []map[string]interface{}{}
	err := model.New(map[string]interface{}{}, nil).Context(ctx).Read(q, &result)
	if err != nil {
		return err
	}
	js, err := json.Marshal(result)
	rsp.Value = string(js)
	return err
}

func (e *Datastore) List(ctx context.Context, req *datastore.ListRequest, rsp *datastore.ListResponse) error {
	log.Info("Received Datastore.List request")
	q := toQuery(req.Query)
	result := []map[string]interface{}{}
	err := model.New(map[string]interface{}{}, nil).Context(ctx).Read(q, &result)
	if err != nil {
		return err
	}
	js, err := json.Marshal(result)
	rsp.Values = string(js)
	return err
}

func (e *Datastore) CreateIndex(ctx context.Context, req *datastore.CreateIndexRequest, rsp *datastore.CreateIndexResponse) error {
	log.Info("Received Datastore.Index request")
	q := toQuery(req.Query)
	return model.New(map[string]interface{}{}, nil).Context(ctx).Delete(q)
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
