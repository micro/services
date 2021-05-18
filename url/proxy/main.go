package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()

		// assuming /u/short-id
		parts := strings.Split(r.URL.Path, "/")
		if len(parts) < 3 {
			return
		}
		// assuming its /u/ for url
		if parts[1] != "u" {
			return
		}

		// get the url id
		//id := parts[2]

		uri := url.URL{
			Scheme: r.URL.Scheme,
			Host:   r.URL.Host,
			Path:   r.URL.Path,
		}

		// call the backend for the url
		rsp, err := http.Get("https://api.m3o.com/url/proxy?shortURL=" + uri.String())
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		defer rsp.Body.Close()

		if rsp.StatusCode != 200 {
			http.Error(w, "unexpected error", 500)
			return
		}

		result := map[string]interface{}{}

		b, err := ioutil.ReadAll(rsp.Body)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		if err := json.Unmarshal(b, &result); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		// get the destination url
		url, _ := result["destinationURL"].(string)
		if len(url) == 0 {
			return
		}

		// return the redirect url to caller
		http.Redirect(w, r, url, 301)
	})

	http.ListenAndServe(":8080", nil)
}
