package handler

import (
	"context"

	"github.com/likexian/doh-go"
	"github.com/likexian/doh-go/dns"

	"github.com/micro/micro/v3/service/errors"
	pb "github.com/micro/services/dns/proto"
)

type Dns struct{}

// Return a new handler
func New() *Dns {
	return &Dns{}
}

func (d *Dns) Query(ctx context.Context, req *pb.QueryRequest, rsp *pb.QueryResponse) error {
	if len(req.Name) == 0 {
		return errors.BadRequest("dns.resolve", "invalid name")
	}
	if len(req.Type) == 0 {
		req.Type = "A"
	}

	c := doh.Use(doh.CloudflareProvider, doh.GoogleProvider)

	// do doh query
	resp, err := c.Query(ctx, dns.Domain(req.Name), dns.Type(req.Type))
	if err != nil {
		return errors.InternalServerError("dns.resolve", err.Error())
	}
	// close the client
	c.Close()

	rsp.Status = int32(resp.Status)
	rsp.TC = resp.TC
	rsp.RD = resp.RD
	rsp.RA = resp.RA
	rsp.AD = resp.AD
	rsp.CD = resp.CD

	for _, q := range resp.Question {
		rsp.Question = append(rsp.Question, &pb.Question{
			Name: q.Name,
			Type: int32(q.Type),
		})
	}

	for _, a := range resp.Answer {
		rsp.Answer = append(rsp.Answer, &pb.Answer{
			Name: a.Name,
			Type: int32(a.Type),
			TTL:  int32(a.TTL),
			Data: a.Data,
		})
	}

	rsp.Provider = resp.Provider

	return nil
}
