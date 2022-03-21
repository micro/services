package tunnel

import (
	"os"

	"net/http"
	"net/http/httputil"
	"net/url"
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

	url, err := url.Parse(r.Host)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// reverse proxy the request
	proxy := httputil.NewSingleHostReverseProxy(url)
	proxy.ServeHTTP(w, r)
}
