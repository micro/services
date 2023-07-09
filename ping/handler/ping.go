package handler

import (
	"bufio"
	"context"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"

	"github.com/go-ping/ping"
	pb "github.com/micro/services/ping/proto"
	"micro.dev/v4/service/errors"
	"micro.dev/v4/service/logger"
)

type Ping struct{}

func (p *Ping) Ip(ctx context.Context, req *pb.IpRequest, rsp *pb.IpResponse) error {
	if len(req.Address) == 0 {
		return errors.BadRequest("ping.ip", "missing address")
	}

	pinger, err := ping.NewPinger(req.Address)
	if err != nil {
		rsp.Status = err.Error()
		return nil
	}

	pinger.Count = 4
	pinger.SetPrivileged(true)
	err = pinger.Run() // Blocks until finished.
	if err != nil {
		rsp.Status = err.Error()
		return nil
	}

	pinger.Stop()
	stats := pinger.Statistics()

	rsp.Status = "OK"
	rsp.Latency = fmt.Sprintf("%v", stats.AvgRtt)

	return nil
}

func (p *Ping) Url(ctx context.Context, req *pb.UrlRequest, rsp *pb.UrlResponse) error {
	if len(req.Address) == 0 {
		return errors.BadRequest("ping.url", "missing address")
	}

	u, err := url.Parse(req.Address)
	if err != nil {
		return errors.BadRequest("ping.url", "failed to parse url: %v", err)
	}

	if u.Scheme == "" {
		u.Scheme = "https"
	}

	if req.Method == "" {
		req.Method = "GET"
	}

	logger.Infof("Calling %s %s", req.Method, u.String())

	hreq, err := http.NewRequest(req.Method, u.String(), nil)
	if err != nil {
		return errors.InternalServerError("ping.url", "failed to make request: %v", err)
	}

	resp, err := http.DefaultClient.Do(hreq)
	if err != nil {
		return errors.InternalServerError("ping.url", "failed to call %v: %v", req.Address, err)
	}
	defer resp.Body.Close()

	// read body
	ioutil.ReadAll(resp.Body)
	// set status
	rsp.Status = resp.Status
	rsp.Code = int32(resp.StatusCode)

	return nil
}

func (p *Ping) Tcp(ctx context.Context, req *pb.TcpRequest, rsp *pb.TcpResponse) error {
	if len(req.Address) == 0 {
		return errors.BadRequest("ping.tcp", "missing address")
	}

	c, err := net.Dial("tcp", req.Address)
	if err != nil {
		rsp.Status = err.Error()
		return nil
	}
	defer c.Close()

	// if no data being sent then just return
	if len(req.Data) == 0 {
		rsp.Status = "OK"
		return nil
	}

	// write data to the connection
	fmt.Fprint(c, req.Data)

	// wait for a response
	data, err := bufio.NewReader(c).ReadString('\n')
	if err != nil {
		rsp.Status = err.Error()
		return nil
	}

	rsp.Status = "OK"
	rsp.Data = data

	return nil
}
