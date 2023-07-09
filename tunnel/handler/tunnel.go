package handler

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"

	pb "github.com/micro/services/tunnel/proto"
	"micro.dev/v4/service/config"
	"micro.dev/v4/service/errors"
	"micro.dev/v4/service/logger"
	"micro.dev/v4/service/store"
)

type Tunnel struct {
	Proxy     string
	Token     string
	Blocklist map[string]bool
}

// loadFile from the blob store
func loadFile(p string) (string, error) {
	name := path.Base(p)

	f, err := os.Create("./" + name)
	if err != nil {
		return "", err
	}
	defer f.Close()

	reader, err := store.DefaultBlobStore.Read(p)
	if err != nil {
		return "", err
	}

	_, err = io.Copy(f, reader)
	return "./" + name, err
}

func New() *Tunnel {
	v, err := config.Get("tunnel.proxy")
	if err != nil {
		logger.Fatalf("failed to get blocklist: %v", err)
	}
	proxy := v.String("")

	v, err = config.Get("tunnel.token")
	if err != nil {
		logger.Fatalf("failed to get token: %v", err)
	}

	token := v.String("")

	v, err = config.Get("tunnel.blocklist")
	if err != nil {
		logger.Fatalf("failed to get blocklist: %v", err)
	}
	path := v.String("")

	// load from blob store if specified
	if strings.HasPrefix(path, "blob://") {
		f, err := loadFile(strings.TrimPrefix(path, "blob://"))
		if err != nil {
			logger.Fatal("failed to load hosts file: %v", err)
		}
		path = f
	}

	f, err := os.ReadFile(path)
	if err != nil {
		logger.Fatal("failed to read hosts file: %v", err)
	}

	blocklist := map[string]bool{}

	for _, host := range strings.Split(string(f), "\n") {
		blocklist[host] = true
	}

	return &Tunnel{
		Blocklist: blocklist,
		Proxy:     proxy,
		Token:     token,
	}
}

func (e *Tunnel) Send(ctx context.Context, req *pb.SendRequest, rsp *pb.SendResponse) error {
	if len(req.Method) == 0 {
		req.Method = "GET"
	}

	// make sure there is a host
	if len(req.Host) == 0 && len(req.Url) == 0 {
		return errors.BadRequest("tunnel.send", "missing host or url")
	}

	var uri *url.URL

	if len(req.Url) > 0 {
		uri, _ = url.Parse(req.Url)
	} else {
		uri = &url.URL{
			Scheme: "https",
			Host:   req.Host,
			Path:   req.Path,
		}
	}

	if len(uri.Scheme) == 0 {
		uri.Scheme = "https"
	}

	vals := url.Values{}
	for k, v := range req.Params {
		vals.Set(k, v)
	}

	if req.Method == "GET" {
		uri.RawQuery = vals.Encode()
	}

	var body io.Reader
	if len(req.Body) > 0 {
		body = bytes.NewReader([]byte(req.Body))
	}

	// check if its a private ip
	if isPrivateIP(uri.Host) {
		logger.Infof("Blocked private host %v", uri.Host)
		return errors.BadRequest("tunnel.send", "cannot send to private ip")
	}

	// check if its in the block list
	if e.Blocklist[strings.ToLower(uri.Host)] {
		logger.Infof("Blocked host %v in blocklist", uri.Host)
		return errors.Forbidden("tunnel.send", "request not allowed")
	}

	// create a new request
	hreq, err := http.NewRequest(req.Method, uri.String(), body)
	if err != nil {
		return errors.BadRequest("tunnel.send", err.Error())
	}

	// set headers
	for k, v := range req.Headers {
		hreq.Header.Set(k, v)
	}

	logger.Infof("Making request %s %s", req.Method, uri.String())

	// set client as default http client
	client := http.DefaultClient

	// use a proxy if specified
	if len(e.Proxy) > 0 {
		hreq.Header.Set("Micro-Endpoint", fmt.Sprintf("%s://%s", uri.Scheme, uri.Host))

		u, _ := url.Parse(e.Proxy)

		// reset the host
		hreq.URL.Scheme = u.Scheme
		hreq.URL.Host = u.Host
		hreq.Host = u.Host
	}

	// set the authorization token
	if len(e.Token) > 0 {
		hreq.Header.Set("Micro-Token", e.Token)
	}

	// make the request
	hrsp, err := client.Do(hreq)
	if err != nil {
		return errors.InternalServerError("tunnel.send", err.Error())
	}
	defer hrsp.Body.Close()

	rsp.Status = hrsp.Status
	rsp.StatusCode = int32(hrsp.StatusCode)
	rsp.Headers = make(map[string]string)
	for k, v := range hrsp.Header {
		rsp.Headers[k] = strings.Join(v, ",")
	}

	b, err := ioutil.ReadAll(hrsp.Body)
	if err != nil {
		return errors.InternalServerError("tunnel.send", "failed to read response")
	}

	// set body
	rsp.Body = string(b)

	return nil
}
