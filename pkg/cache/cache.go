// Cache provides a simple marshaling layer on top of the store
package cache

import (
	"encoding/json"
	"time"

	"github.com/micro/micro/v3/service/store"
)

func Get(key string, val interface{}) error {
	recs, err := store.Read(key, store.ReadLimit(1))
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

func Put(key string, val interface{}, expires time.Time) error {
	b, err := json.Marshal(val)
	if err != nil {
		return err
	}
	expiry := expires.Sub(time.Now())
	if expiry < time.Duration(0) {
		expiry = time.Duration(0)
	}
	return store.Write(&store.Record{
		Key:    key,
		Value:  b,
		Expiry: expiry,
	})
}

func Delete(key string) error {
	return store.Delete(key)
}
