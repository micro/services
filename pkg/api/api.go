// Package api is a helper for calling external third party apis
package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func Get(url string, rsp interface{}) error {
	resp, err := http.Get(url)
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
