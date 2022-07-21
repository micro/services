package proxy

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

var (
	// API Key
	APIKey = os.Getenv("MICRO_API_KEY")

	// API Url
	APIHost = "https://api.m3o.com"

	// host to proxy for
	Host = "m3o.one"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	uri := url.URL{
		Scheme: r.URL.Scheme,
		Host:   r.URL.Host,
		Path:   r.URL.Path,
	}

	// if the host is blank we have to set it
	if len(uri.Host) == 0 {
		if len(r.Host) > 0 {
			log.Printf("[url/proxy] Host is set from r.Host %v", r.Host)
			uri.Host = r.Host
		} else {
			log.Printf("[url/proxy] Host is nil, defaulting to: %v", Host)
			uri.Host = Host
			uri.Scheme = "https"
		}
	}

	if len(uri.Scheme) == 0 {
		uri.Scheme = "https"
	}

	apiURL := APIHost + "/url/resolve"

	// use /v1/
	if len(APIKey) > 0 {
		apiURL = APIHost + "/v1/url/resolve"
	}

	// make new request
	log.Printf("[url/proxy] Calling: %v", apiURL+"?shortURL="+uri.String())
	req, err := http.NewRequest("GET", apiURL+"?shortURL="+uri.String(), nil)
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
		log.Printf("[url/proxy] Error calling api: status: %v %v", rsp.StatusCode, string(b))
		http.Error(w, "unexpected error", 500)
		return
	}

	result := map[string]interface{}{}

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
	http.Redirect(w, r, url, 302)
}
