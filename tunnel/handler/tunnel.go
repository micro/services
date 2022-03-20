package handler

import (
	"bytes"
	"context"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/micro/micro/v3/service/errors"
	"github.com/micro/micro/v3/service/logger"
	pb "github.com/micro/services/tunnel/proto"
)

type Tunnel struct{}

func New() *Tunnel {
	return &Tunnel{}
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

	if isPrivateIP(uri.Host) {
		return errors.BadRequest("tunnel.send", "cannot send to private ip")
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

	// make the request
	hrsp, err := http.DefaultClient.Do(hreq)
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
