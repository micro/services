// package model helps with data modelling on top of the store
package model

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/micro/micro/v3/service/store"
)

var (
	ErrNotFound      = errors.New("not found")
	ErrAlreadyExists = errors.New("already exists")
)

type Entity interface {
	// The primary key
	Key(ctx context.Context) string
	// The index for the entity
	Index(ctx context.Context) string
	// The raw value of the entity
	Value() interface{}
}

type Query struct {
	Limit  uint
	Offset uint
	Order  string
}

func Create(ctx context.Context, e Entity) error {
	key := e.Key(ctx)
	val := e.Value()
	idx := e.Index(ctx)

	// read the existing record
	recs, err := store.Read(key, store.ReadLimit(1))
	if err != nil && err != store.ErrNotFound {
		return err
	}

	if len(recs) > 0 {
		return ErrAlreadyExists
	}

	// write the record
	if err := store.Write(store.NewRecord(key, val)); err != nil {
		return err
	}

	// only write the index if it exists
	if len(idx) == 0 {
		return nil
	}

	// write the index
	return store.Write(store.NewRecord(idx, val))
}

func ReadIndex(ctx context.Context, e Entity) error {
	recs, err := store.Read(e.Index(ctx), store.ReadLimit(1))
	if err == store.ErrNotFound {
		return ErrNotFound
	} else if err != nil {
		return err
	}
	if len(recs) == 0 {
		return ErrNotFound
	}
	return recs[0].Decode(e)
}

func Read(ctx context.Context, e Entity) error {
	recs, err := store.Read(e.Key(ctx), store.ReadLimit(1))
	if err == store.ErrNotFound {
		return ErrNotFound
	} else if err != nil {
		return err
	}
	if len(recs) == 0 {
		return ErrNotFound
	}
	return recs[0].Decode(e)
}

func Update(ctx context.Context, e Entity) error {
	key := e.Key(ctx)
	val := e.Value()
	idx := e.Index(ctx)

	// write the record
	if err := store.Write(store.NewRecord(key, val)); err != nil {
		return err
	}

	// only write the index if it exists
	if len(idx) == 0 {
		return nil
	}

	// write the index
	return store.Write(store.NewRecord(idx, val))
}

func List(ctx context.Context, e Entity, rsp interface{}, q Query) error {
	opts := []store.ReadOption{
		store.ReadPrefix(),
	}

	if q.Limit > 0 {
		opts = append(opts, store.ReadLimit(q.Limit))
	}
	if q.Offset > 0 {
		opts = append(opts, store.ReadOffset(q.Offset))
	}
	if len(q.Order) > 0 {
		if q.Order == "desc" {
			opts = append(opts, store.ReadOrder(store.OrderDesc))
		} else {
			opts = append(opts, store.ReadOrder(store.OrderAsc))
		}
	}

	recs, err := store.Read(e.Index(ctx), opts...)
	if err != nil {
		return err
	}

	jsBuffer := []byte("[")

	for i, rec := range recs {
		jsBuffer = append(jsBuffer, rec.Value...)
		if i < len(recs)-1 {
			jsBuffer = append(jsBuffer, []byte(",")...)
		}
	}

	jsBuffer = append(jsBuffer, []byte("]")...)
	return json.Unmarshal(jsBuffer, rsp)
}

func Delete(ctx context.Context, e Entity) error {
	key := e.Key(ctx)
	idx := e.Index(ctx)

	if len(key) > 0 {
		if err := store.Delete(key); err != nil && err != store.ErrNotFound {
			return err
		}
	}

	recs, err := store.Read(idx, store.ReadPrefix())
	if err != nil && err != store.ErrNotFound {
		return err
	}

	// delete every record by index
	for _, rec := range recs {
		var val interface{}
		if err := rec.Decode(val); err != nil || val == nil {
			continue
		}
		// convert to an entity
		e, ok := val.(Entity)
		if !ok {
			continue
		}
		if err := store.Delete(e.Key(ctx)); err != store.ErrNotFound {
			return err
		}
	}

	return nil
}
