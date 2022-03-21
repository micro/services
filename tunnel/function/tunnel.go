package tunnel

import (
	"os"
	"log"
	"net/http"
	"net/http/httputil"
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

	// reverse proxy the request
	log.Printf("Proxying request to: %v", r.URL.Host)
	proxy := httputil.NewSingleHostReverseProxy(r.URL)
	proxy.ServeHTTP(w, r)
}
