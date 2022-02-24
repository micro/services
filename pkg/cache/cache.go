// Cache provides a simple marshaling layer on top of the store
package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/hashicorp/golang-lru"
	"github.com/micro/micro/v3/service/logger"
	"github.com/micro/micro/v3/service/store"
	"github.com/micro/services/pkg/tenant"
	"github.com/peterbourgon/diskv/v3"
)

type Cache interface {
	// Context returns a tenant scoped Cache
	Context(ctx context.Context) Cache
	Get(key string, val interface{}) (time.Time, error)
	Set(key string, val interface{}, expires time.Time) error
	Delete(key string) error
	Increment(key string, val int64) (int64, error)
	Decrement(key string, val int64) (int64, error)
	ListKeys() ([]string, error)
	Close() error
}

type cache struct {
	sync.RWMutex
	closed chan bool
	LRU    *lru.Cache
	Disk   *diskv.Diskv
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
	d := diskv.New(diskv.Options{
		BasePath: "cache",
	})

	return &cache{
		LRU:   l,
		Disk:  d,
		Store: st,
	}
}

func (c *cache) run() {
	for {
		select {
		case <-c.closed:
			return
		case <-time.After(time.Hour):
			c.Disk.EraseAll()
		}
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
		closed: make(chan bool),
		LRU:    c.LRU,
		Disk:   c.Disk,
		Store:  c.Store,
		Prefix: t,
	}
}

func (c *cache) Close() error {
	c.Lock()
	defer c.Unlock()

	select {
	case <-c.closed:
		return nil
	default:
		close(c.closed)
	}

	return nil
}

func (c *cache) Get(key string, val interface{}) (time.Time, error) {
	k := c.Key(key)

	// try the LRU
	v, ok := c.LRU.Get(k)
	if ok {
		i := v.(*item)

		// check if the item expired
		if !i.expires.IsZero() && i.expires.Sub(time.Now()).Seconds() < 0 {
			// remove it
			c.LRU.Remove(k)
			return time.Time{}, ErrNotFound
		}

		// otherwise unmarshal and return it
		return i.expires, json.Unmarshal(i.val, val)
	}

	logger.Infof("Cache miss for %v", k)

	if c.Disk == nil {
		c.Disk = diskv.New(diskv.Options{
			BasePath: "cache",
		})
	}

	// read from disk
	b, err := c.Disk.Read(k)
	if err == nil && len(b) > 0 {
		var i item
		if err := json.Unmarshal(b, &i); err == nil {
			if !i.expires.IsZero() && i.expires.Sub(time.Now()).Seconds() < 0 {
				c.Disk.Erase(k)
				return time.Time{}, ErrNotFound
			}
			return i.expires, json.Unmarshal(i.val, val)
		}
	}

	// otherwise check  the store
	if c.Store == nil {
		c.Store = store.DefaultStore
	}

	recs, err := c.Store.Read(k, store.ReadLimit(1))
	if err != nil && err == store.ErrNotFound {
		return time.Time{}, ErrNotFound
	} else if err != nil {
		return time.Time{}, err
	}
	if len(recs) == 0 {
		return time.Time{}, ErrNotFound
	}
	if err := json.Unmarshal(recs[0].Value, val); err != nil {
		return time.Time{}, err
	}

	// put it in the cache for future use
	// set in the lru
	rec := recs[0]
	expires := time.Time{}
	if rec.Expiry > time.Duration(0) {
		expires = time.Now().Add(rec.Expiry)
	}

	vi := &item{key: rec.Key, val: rec.Value, expires: expires}
	c.LRU.Add(rec.Key, vi)

	b, _ = json.Marshal(vi)
	// put on disk
	c.Disk.Write(rec.Key, b)

	return vi.expires, nil
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
	vi := &item{key: rec.Key, val: rec.Value, expires: expires}
	c.LRU.Add(rec.Key, vi)
	b, _ = json.Marshal(vi)
	// put on disk
	c.Disk.Write(rec.Key, b)
	return nil
}

func (c *cache) Delete(key string) error {
	if c.Store == nil {
		c.Store = store.DefaultStore
	}

	k := c.Key(key)
	// remove from the lru
	c.LRU.Remove(k)
	// remove from disk
	c.Disk.Erase(k)
	// delete from the store
	return c.Store.Delete(k)
}

func (c *cache) Increment(key string, value int64) (int64, error) {
	c.Lock()
	defer c.Unlock()

	var val interface{}
	if _, err := c.Get(key, &val); err != nil && err != ErrNotFound {
		return 0, err
	}

	var counter int64

	switch val.(type) {
	case string:
		// try to convert to number
		a, err := strconv.ParseInt(val.(string), 10, 64)
		if err != nil {
			return 0, err
		}
		counter = a
	case int64:
		counter = val.(int64)
	case int32:
		counter = int64(val.(int32))
	case int:
		counter = int64(val.(int))
	case nil:
		counter = 0
	default:
		return 0, errors.New("value is not an integer")
	}

	counter += value

	if err := c.Set(key, fmt.Sprintf("%v",counter), time.Time{}); err != nil {
		return counter, err
	}
	return counter, nil
}

func (c *cache) Decrement(key string, value int64) (int64, error) {
	c.Lock()
	defer c.Unlock()

	var val interface{}
	if _, err := c.Get(key, &val); err != nil && err != ErrNotFound {
		return 0, err
	}

	var counter int64

	switch val.(type) {
	case string:
		// try to convert to number
		a, err := strconv.ParseInt(val.(string), 10, 64)
		if err != nil {
			return 0, err
		}
		counter = a
	case int64:
		counter = val.(int64)
	case int32:
		counter = int64(val.(int32))
	case int:
		counter = int64(val.(int))
	case nil:
		counter = 0
	default:
		return 0, errors.New("value is not an integer")
	}

	counter -= value

	if err := c.Set(key, fmt.Sprintf("%v", counter), time.Time{}); err != nil {
		return counter, err
	}
	return counter, nil
}

func (c *cache) ListKeys() ([]string, error) {
	c.RLock()
	defer c.RUnlock()

	if c.Store == nil {
		c.Store = store.DefaultStore
	}

	prefix := c.Key("")

	recKeys, err := c.Store.List(store.ListPrefix(prefix))
	if err != nil {
		return nil, err
	}

	var keys []string

	for _, key := range recKeys {
		keys = append(keys, strings.TrimPrefix(key, prefix))
	}

	return keys, nil
}

func Context(ctx context.Context) Cache {
	return DefaultCache.Context(ctx)
}

func Get(key string, val interface{}) (time.Time, error) {
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

func ListKeys() ([]string, error) {
	return DefaultCache.ListKeys()
}
