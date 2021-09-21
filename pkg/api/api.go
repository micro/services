// Package api is a helper for calling external third party apis
package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

var (
	keys = map[string]string{}
)

// Set a key within the header
func SetKey(k, v string) {
	keys[k] = v
}

// Get a url and unmarshal a json body into the given value
func Get(url string, rsp interface{}) error {
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

	if resp.StatusCode != 200 {
		return fmt.Errorf("Non 200 response %v: %v", resp.StatusCode, string(b))
	}

	return json.Unmarshal(b, rsp)
}
