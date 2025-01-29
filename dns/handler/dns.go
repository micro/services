package handler

import (
	"context"

	"github.com/likexian/doh-go"
	"github.com/likexian/doh-go/dns"
	"github.com/likexian/whois"
	"github.com/likexian/whois-parser"
	"github.com/micro/micro/v5/service/errors"
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

func (d *Dns) Whois(ctx context.Context, req *pb.WhoisRequest, rsp *pb.WhoisResponse) error {
	if len(req.Domain) == 0 {
		return errors.BadRequest("dns.whois", "missing domain")
	}

	result, err := whois.Whois(req.Domain)
	if err != nil {
		return err
	}

	// parse the response
	res, err := whoisparser.Parse(result)
	if err != nil {
		return err
	}

	rsp.Id = res.Domain.ID
	rsp.Domain = res.Domain.Domain
	rsp.WhoisServer = res.Domain.WhoisServer
	rsp.RegistrarUrl = res.Registrar.ReferralURL
	rsp.Created = res.Domain.CreatedDate
	rsp.Updated = res.Domain.UpdatedDate
	rsp.Expiry = res.Domain.ExpirationDate
	rsp.Status = res.Domain.Status
	rsp.Registrar = res.Registrar.Name
	rsp.RegistrarId = res.Registrar.ID
	rsp.AbuseEmail = res.Registrar.Email
	rsp.AbusePhone = res.Registrar.Phone
	rsp.Nameservers = res.Domain.NameServers
	return nil
}
