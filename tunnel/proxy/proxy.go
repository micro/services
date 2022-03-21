package proxy

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
)

var (
	token = os.Getenv("TOKEN")
)

func Send(w http.ResponseWriter, r *http.Request) {
	// check the auth token
	x := r.Header.Get("Micro-Token")
	if x != token {
		http.Error(w, "unauthorized", 401)
		return
	}

	// don't forward the token
	r.Header.Del("Micro-Token")

	// get the endpoint
	u := r.Header.Get("Micro-Endpoint")
	if len(u) == 0 {
		return
	}

	// delete the endpoint header
	r.Header.Del("Micro-Endpoint")

	// parse the request url
	uri, _ := url.Parse(u)

	r.Host = uri.Host
	r.URL.Host = uri.Host
	r.URL.Scheme = uri.Scheme

	// reverse proxy the request
	fmt.Printf("Proxying request to: %v", uri.String())

	// proxy the request
	proxy := httputil.NewSingleHostReverseProxy(r.URL)
	proxy.ServeHTTP(w, r)
}
