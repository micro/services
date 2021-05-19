// Cache provides a simple marshaling layer on top of the store
package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/micro/micro/v3/service/logger"
	"github.com/hashicorp/golang-lru"
	"github.com/micro/micro/v3/service/store"
	"github.com/micro/services/pkg/tenant"
)

type Cache interface {
	// Context returns a tenant scoped Cache
	Context(ctx context.Context) Cache
	Get(key string, val interface{}) error
	Set(key string, val interface{}, expires time.Time) error
	Delete(key string) error
	Increment(key string, val int64) (int64, error)
	Decrement(key string, val int64) (int64, error)
}

type cache struct {
	sync.Mutex
	LRU    *lru.Cache
	Store  store.Store
	Prefix string
}

type item struct {
	key     string
	val     []byte
	expires time.Time
}

var (
	DefaultCacheSize = 1000

	DefaultCache = New(nil)

	ErrNotFound = errors.New("not found")
)

func New(st store.Store) Cache {
	l, _ := lru.New(DefaultCacheSize)
	return &cache{
		LRU:   l,
		Store: st,
	}
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
		LRU:    c.LRU,
		Store:  c.Store,
		Prefix: t,
	}
}

func (c *cache) Get(key string, val interface{}) error {
	k := c.Key(key)

	// try the LRU
	v, ok := c.LRU.Get(k)
	if ok {
		i := v.(*item)

		// check if the item expired
		if !i.expires.IsZero() && i.expires.Sub(time.Now()).Seconds() < 0 {
			// remove it
			c.LRU.Remove(k)
			return ErrNotFound
		}

		// otherwise unmarshal and return it
		return json.Unmarshal(i.val, val)
	}

	logger.Infof("Cache miss for %v", k)

	// otherwise check  the store
	if c.Store == nil {
		c.Store = store.DefaultStore
	}

	recs, err := c.Store.Read(k, store.ReadLimit(1))
	if err != nil && err == store.ErrNotFound {
		return ErrNotFound
	} else if err != nil {
	}
	if len(recs) == 0 {
		return ErrNotFound
	}
	if err := json.Unmarshal(recs[0].Value, val); err != nil {
		return err
	}

	// put it in the cache for future use
	// set in the lru
	rec := recs[0]
	expires := time.Time{}
	if rec.Expiry > time.Duration(0) {
		expires = time.Now().Add(rec.Expiry)
	}
	c.LRU.Add(rec.Key, &item{key: rec.Key, val: rec.Value, expires: expires})

	return nil
}

func (c *cache) Set(key string, val interface{}, expires time.Time) error {
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
	rec := &store.Record{
		Key:    c.Key(key),
		Value:  b,
		Expiry: expiry,
	}
	if err := c.Store.Write(rec); err != nil {
		return err
	}

	// set in the lru
	c.LRU.Add(rec.Key, &item{key: rec.Key, val: rec.Value, expires: expires})
	return nil
}

func (c *cache) Delete(key string) error {
	if c.Store == nil {
		c.Store = store.DefaultStore
	}

	k := c.Key(key)
	// remove from the lru
	c.LRU.Remove(k)
	// delete from the store
	return c.Store.Delete(k)
}

func (c *cache) Increment(key string, value int64) (int64, error) {
	c.Lock()
	defer c.Unlock()

	var val int64
	if err := c.Get(key, &val); err != nil && err != ErrNotFound {
		return 0, err
	}
	val += value
	if err := c.Set(key, val, time.Time{}); err != nil {
		return val, err
	}
	return val, nil
}

func (c *cache) Decrement(key string, value int64) (int64, error) {
	c.Lock()
	defer c.Unlock()

	var val int64
	if err := c.Get(key, &val); err != nil && err != ErrNotFound {
		return 0, err
	}
	val -= value
	if err := c.Set(key, val, time.Time{}); err != nil {
		return val, err
	}
	return val, nil
}

func Context(ctx context.Context) Cache {
	return DefaultCache.Context(ctx)
}

func Get(key string, val interface{}) error {
	return DefaultCache.Get(key, val)
}

func Set(key string, val interface{}, expires time.Time) error {
	return DefaultCache.Set(key, val, expires)
}

func Delete(key string) error {
	return DefaultCache.Delete(key)
}

func Increment(key string, val int64) (int64, error) {
	return DefaultCache.Increment(key, val)
}

func Decrement(key string, val int64) (int64, error) {
	return DefaultCache.Decrement(key, val)
}
