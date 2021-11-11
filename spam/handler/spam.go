package handler

import (
	"context"
	"strings"

	"github.com/Teamwork/spamc"
	"github.com/micro/micro/v3/service/config"
	"github.com/micro/micro/v3/service/errors"
	log "github.com/micro/micro/v3/service/logger"
	spam "github.com/micro/services/spam/proto"
)

type conf struct {
	SpamdAddress string `json:"spamd_address"`
}

type Spam struct {
	c      conf
	client *spamc.Client
}

func New() *Spam {
	val, err := config.Get("micro.spam")
	if err != nil {
		log.Fatalf("Failed to load config")
	}
	c := conf{}
	if err := val.Scan(&c); err != nil {
		log.Fatalf("Failed to load config")
	}
	return &Spam{
		c:      c,
		client: spamc.New(c.SpamdAddress, nil),
	}
}

func (s *Spam) Check(ctx context.Context, request *spam.CheckRequest, response *spam.CheckResponse) error {
	hdr := spamc.Header{}
	for k, v := range request.Headers {
		hdr.Set(k, v)
	}
	rc, err := s.client.Report(ctx, strings.NewReader(request.EmailBody), hdr)
	if err != nil {
		log.Errorf("Error checking spamd %s", err)
		return errors.InternalServerError("spam.Check", "Error checking spam")
	}
	response.IsSpam = rc.IsSpam
	response.Score = rc.Score
	response.Report = rc.Report.String()
	return nil
}
