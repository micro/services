// Cache provides a simple marshaling layer on top of the store
package cache

import (
	"encoding/json"
	"time"

	"github.com/micro/micro/v3/service/store"
)

type Cache interface {
	Get(key string, val interface{}) error
	Put(key string, val interface{}, expires time.Time) error
	Delete(key string) error
}

type cache struct {
	Store store.Store
}

var (
	DefaultCache = New(nil)
)

func New(st store.Store) Cache {
	return &cache{st}
}

func (c *cache) Get(key string, val interface{}) error {
	if c.Store == nil {
		c.Store = store.DefaultStore
	}

	recs, err := c.Store.Read(key, store.ReadLimit(1))
	if err != nil {
		return err
	}
	if len(recs) == 0 {
		return store.ErrNotFound
	}
	if err := json.Unmarshal(recs[0].Value, val); err != nil {
		return err
	}
	return nil
}

func (c *cache) Put(key string, val interface{}, expires time.Time) error {
	if c.Store == nil {
		c.Store = store.DefaultStore
	}
	b, err := json.Marshal(val)
	if err != nil {
		return err
	}
	expiry := expires.Sub(time.Now())
	if expiry < time.Duration(0) {
		expiry = time.Duration(0)
	}
	return c.Store.Write(&store.Record{
		Key:    key,
		Value:  b,
		Expiry: expiry,
	})
}

func (c *cache) Delete(key string) error {
	if c.Store == nil {
		c.Store = store.DefaultStore
	}
	return c.Store.Delete(key)
}

func Get(key string, val interface{}) error {
	return DefaultCache.Get(key, val)
}

func Put(key string, val interface{}, expires time.Time) error {
	return DefaultCache.Put(key, val, expires)
}

func Delete(key string) error {
	return DefaultCache.Delete(key)
}
