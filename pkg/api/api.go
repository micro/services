// Package api is a helper for calling external third party apis
package api

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"path"
	"sync"
	"time"
)

var (
	keys = map[string]string{}

	mtx      sync.RWMutex
	cache    map[string]interface{}
	cacheTTL time.Duration
)

// Set a key within the header
func SetKey(k, v string) {
	keys[k] = v
}

// Set the cache
func SetCache(v bool, ttl time.Duration) {
	cache = make(map[string]interface{})
	cacheTTL = ttl
}

// Get a url and unmarshal a json body into the given value
func Get(url string, rsp interface{}) error {
	// check the cache if its enabled
	mtx.RLock()
	if cache != nil {
		if val, ok := cache[url]; ok {
			mtx.RUnlock()
			return json.Unmarshal(val.([]byte), rsp)
		}
	}
	mtx.RUnlock()

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	for k, v := range keys {
		req.Header.Set(k, v)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode >= 400 {
		return fmt.Errorf("Non 200 response %v: %v", resp.StatusCode, string(b))
	}

	if cache != nil {
		mtx.Lock()
		// cache the value
		cache[url] = b

		// delete it when the ttl expires
		if cacheTTL > time.Duration(0) {
			go func() {
				time.Sleep(cacheTTL)
				mtx.Lock()
				delete(cache, url)
				mtx.Unlock()
			}()
		}
		mtx.Unlock()
	}

	return json.Unmarshal(b, rsp)
}

// Post a url and unmarshal a json body into the given value
func Post(url string, ureq, rsp interface{}) error {
	b, err := json.Marshal(ureq)
	if err != nil {
		return err
	}

	// encode the data
	data := base64.StdEncoding.EncodeToString(b)

	// create a key
	key := path.Join(url, data)

	// check the cache if its enabled
	mtx.RLock()
	if cache != nil {
		if val, ok := cache[key]; ok {
			mtx.RUnlock()
			return json.Unmarshal(val.([]byte), rsp)
		}
	}
	mtx.RUnlock()

	req, err := http.NewRequest("POST", url, bytes.NewReader(b))
	if err != nil {
		return err
	}

	for k, v := range keys {
		req.Header.Set(k, v)
	}

	if v := req.Header.Get("Content-Type"); len(v) == 0 {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	b, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode >= 400 {
		return fmt.Errorf("Non 200 response %v: %v", resp.StatusCode, string(b))
	}

	if cache != nil {
		mtx.Lock()

		// cache the value
		cache[key] = b

		// delete it when the ttl expires
		if cacheTTL > time.Duration(0) {
			go func() {
				time.Sleep(cacheTTL)
				mtx.Lock()
				delete(cache, key)
				mtx.Unlock()
			}()
		}
		mtx.Unlock()
	}

	return json.Unmarshal(b, rsp)
}

// Call the DELETE method for a URL
func Delete(url string, rsp interface{}) error {
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return err
	}

	for k, v := range keys {
		req.Header.Set(k, v)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode >= 400 {
		return fmt.Errorf("Non 200 response %v: %v", resp.StatusCode, string(b))
	}

	if cache != nil {
		mtx.Lock()
		// cache the value
		delete(cache, url)
		mtx.Unlock()
	}

	return json.Unmarshal(b, rsp)
}
