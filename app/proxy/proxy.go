package proxy

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"
)

var (
	// API Key
	APIKey = os.Getenv("MICRO_API_KEY")

	// API Url
	APIHost = "https://api.m3o.com"

	// host to proxy for Apps
	AppHost = "m3o.app"
)

var (
	mtx sync.RWMutex

	// local cache
	appMap = map[string]*backend{}
)

type backend struct {
	url     *url.URL
	created time.Time
}

func Handler(w http.ResponseWriter, r *http.Request) {
	// no subdomain
	if r.Host == AppHost {
		return
	}

	// lookup the app map
	mtx.RLock()
	bk, ok := appMap[r.Host]
	mtx.RUnlock()

	// check the url map
	if ok && time.Since(bk.created) < time.Minute {
		r.Host = bk.url.Host
		r.Header.Set("Host", r.Host)
		httputil.NewSingleHostReverseProxy(bk.url).ServeHTTP(w, r)
		return
	}

	subdomain := strings.TrimSuffix(r.Host, "."+AppHost)

	// only process one part for now
	parts := strings.Split(subdomain, ".")
	if len(parts) > 1 {
		log.Print("[app/proxy] more parts than expected", parts)
		return
	}

	// currently service id is the subdomain
	id := subdomain

	log.Printf("[app/proxy] resolving host %s to id %s\n", r.Host, id)

	apiURL := APIHost + "/app/resolve"

	// use /v1/
	if len(APIKey) > 0 {
		apiURL = APIHost + "/v1/app/resolve"
	}

	// make new request
	log.Printf("[app/proxy] Calling: %v", apiURL+"?id="+id)
	req, err := http.NewRequest("GET", apiURL+"?id="+id, nil)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	if req.Header == nil {
		req.Header = make(http.Header)
	}

	// set the api key after we're given the header
	if len(APIKey) > 0 {
		req.Header.Set("Authorization", "Bearer "+APIKey)
	}

	// call the backend for the url
	rsp, err := http.DefaultClient.Do(req)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	defer rsp.Body.Close()

	b, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	if rsp.StatusCode != 200 {
		log.Printf("[app/proxy] Error calling api: status: %v %v", rsp.StatusCode, string(b))
		http.Error(w, "unexpected error", 500)
		return
	}

	result := map[string]interface{}{}

	if err := json.Unmarshal(b, &result); err != nil {
		log.Print("[app/proxy] failed to unmarshal response")
		http.Error(w, err.Error(), 500)
		return
	}

	// get the destination url
	u, _ := result["url"].(string)
	if len(u) == 0 {
		log.Print("[app/proxy] no response url")
		return
	}

	uri, err := url.Parse(u)
	if err != nil {
		log.Print("[app/proxy] failed to parse url", err.Error())
		return
	}

	mtx.Lock()
	appMap[r.Host] = &backend{
		url:     uri,
		created: time.Now(),
	}
	mtx.Unlock()

	r.Host = uri.Host
	r.Header.Set("Host", r.Host)

	httputil.NewSingleHostReverseProxy(uri).ServeHTTP(w, r)
}
