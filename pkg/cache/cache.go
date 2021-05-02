// Cache provides a simple marshaling layer on top of the store
package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/micro/micro/v3/service/store"
	"github.com/micro/services/pkg/tenant"
)

type Cache interface {
	// Context returns a tenant scoped Cache
	Context(ctx context.Context) Cache
	Get(key string, val interface{}) error
	Put(key string, val interface{}, expires time.Time) error
	Delete(key string) error
}

type cache struct {
	Store  store.Store
	Prefix string
}

var (
	DefaultCache = New(nil)
)

func New(st store.Store) Cache {
	return &cache{Store: st}
}

func (c *cache) Key(k string) string {
	if len(c.Prefix) > 0 {
		return fmt.Sprintf("%s/%s", c.Prefix, k)
	}
	return k
}

func (c *cache) Context(ctx context.Context) Cache {
	t, ok := tenant.FromContext(ctx)
	if !ok {
		return c
	}
	return &cache{
		Store:  c.Store,
		Prefix: t,
	}
}

func (c *cache) Get(key string, val interface{}) error {
	if c.Store == nil {
		c.Store = store.DefaultStore
	}

	recs, err := c.Store.Read(c.Key(key), store.ReadLimit(1))
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
		Key:    c.Key(key),
		Value:  b,
		Expiry: expiry,
	})
}

func (c *cache) Delete(key string) error {
	if c.Store == nil {
		c.Store = store.DefaultStore
	}
	return c.Store.Delete(c.Key(key))
}

func Context(ctx context.Context) Cache {
	return DefaultCache.Context(ctx)
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
