package handler

import (
	"bytes"
	"context"
	"fmt"
	"time"

	"github.com/Teamwork/spamc"
	"github.com/micro/micro/v3/service/config"
	"github.com/micro/micro/v3/service/errors"
	log "github.com/micro/micro/v3/service/logger"
	spam "github.com/micro/services/spam/proto"
	"gopkg.in/gomail.v2"
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
	if len(request.EmailBody) == 0 && len(request.TextBody) == 0 && len(request.HtmlBody) == 0 {
		return errors.BadRequest("spam.Classify", "Missing one of email_body, html_body, text_body")
	}
	bf := bytes.Buffer{}

	if len(request.EmailBody) > 0 {
		bf.WriteString(request.EmailBody)
	} else {
		m := gomail.NewMessage()

		if len(request.To) > 0 {
			m.SetHeader("To", request.To)
		}
		if len(request.From) > 0 {
			m.SetHeader("From", request.From)
		}
		if len(request.Subject) > 0 {
			m.SetHeader("Subject", request.Subject)
		}
		m.SetHeader("Date", time.Now().Format(time.RFC1123Z))
		if len(request.TextBody) > 0 {
			m.SetBody("text/plain", request.TextBody)
		}
		if len(request.HtmlBody) > 0 {
			m.SetBody("text/html", request.HtmlBody)
		}
		if _, err := m.WriteTo(&bf); err != nil {
			log.Errorf("Error classifying email %s", err)
			return errors.InternalServerError("spam.Classify", "Error classifying email")
		}

	}
	rc, err := s.client.Report(ctx, &bf, spamc.Header{}.Set("Content-Length", fmt.Sprintf("%d", bf.Len())))
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
