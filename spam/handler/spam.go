package handler

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"net/textproto"
	"time"

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

func (s *Spam) Classify(ctx context.Context, request *spam.ClassifyRequest, response *spam.ClassifyResponse) error {
	if len(request.EmailBody) == 0 {
		return errors.BadRequest("spam.Classify", "missing email_body")
	}

	bf := bytes.NewBufferString("")
	tp := textproto.NewWriter(bufio.NewWriter(bf))

	if len(request.To) > 0 {
		tp.PrintfLine("To: %v", request.To)
	}
	if len(request.From) > 0 {
		tp.PrintfLine("From: %v", request.From)
	}
	if len(request.Subject) > 0 {
		tp.PrintfLine("Subject: %v", request.Subject)
	}
	tp.PrintfLine("Date: %s", time.Now().Format(time.RFC1123Z))
	tp.PrintfLine("")
	tp.PrintfLine("%v", request.EmailBody)
	rc, err := s.client.Report(ctx, bf, spamc.Header{}.Set("Content-Length", fmt.Sprintf("%d", bf.Len())))
	if err != nil {
		log.Errorf("Error checking spamd %s", err)
		return errors.InternalServerError("spam.Classify", "Error classifying email")
	}
	response.IsSpam = rc.IsSpam
	response.Score = rc.Score

	response.Details = []string{}
	for _, v := range rc.Report.Table {
		response.Details = append(response.Details, fmt.Sprintf("%s, %s, %v", v.Rule, v.Description, v.Points))
	}
	return nil
}
