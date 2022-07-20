package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/mail"
	"regexp"

	"github.com/micro/micro/v3/service"
	"github.com/micro/micro/v3/service/client"
	"github.com/micro/micro/v3/service/config"
	"github.com/micro/micro/v3/service/errors"
	log "github.com/micro/micro/v3/service/logger"
	"github.com/micro/micro/v3/service/store"
	pb "github.com/micro/services/email/proto"
	"github.com/micro/services/pkg/tenant"
	spampb "github.com/micro/services/spam/proto"
)

const (
	prefixUserID     = "byUserID"
	prefixSendgridID = "bySendgridID"
)

type Sent struct {
	UserID        string
	SendgridMsgID string
}

type sendgridConf struct {
	Key       string `json:"key"`
	EmailFrom string `json:"email_from"`
	PoolName  string `json:"ip_pool_name"`
}

func NewEmailHandler(svc *service.Service) *Email {
	c := sendgridConf{}
	val, err := config.Get("sendgridapi")
	if err != nil {
		log.Warnf("Error getting config: %v", err)
	}
	err = val.Scan(&c)
	if err != nil {
		log.Warnf("Error scanning config: %v", err)
	}
	if len(c.Key) == 0 {
		log.Fatalf("Sendgrid API key not configured")
	}
	return &Email{
		c,
		spampb.NewSpamService("spam", svc.Client()),
	}
}

type Email struct {
	config  sendgridConf
	spamSvc spampb.SpamService
}

// validEmail does very light validation
func validEmail(email string) bool {
	if len(email) == 0 {
		return false
	}
	m, err := regexp.MatchString("^\\S+@\\S+$", email)
	if err != nil {
		return false
	}
	return m
}

func (e *Email) Send(ctx context.Context, request *pb.SendRequest, response *pb.SendResponse) error {
	if len(request.From) == 0 {
		return errors.BadRequest("email.send.validation", "Missing from name")
	}
	if !validEmail(request.To) {
		return errors.BadRequest("email.send.validation", "Invalid to address")
	}
	if len(request.Subject) == 0 {
		return errors.BadRequest("email.send.validation", "Missing subject")
	}
	if len(request.TextBody) == 0 && len(request.HtmlBody) == 0 {
		return errors.BadRequest("email.send.validation", "Missing email body")
	}

	spamReq := &spampb.ClassifyRequest{
		TextBody: request.TextBody,
		HtmlBody: request.HtmlBody,
		To:       request.To,
		From:     request.From,
		Subject:  request.Subject,
	}
	rsp, err := e.spamSvc.Classify(ctx, spamReq, client.WithAuthToken())
	if err != nil || rsp.IsSpam {
		log.Errorf("Error validating email %s %v", err, rsp)
		return errors.InternalServerError("email.send", "Error validating email")
	}

	if err := e.sendEmail(ctx, request); err != nil {
		log.Errorf("Error sending email: %v\n", err)
		return errors.InternalServerError("email.sendemail", "Error sending email")
	}

	return nil
}

// sendEmail sends an email via the sendgrid API
// Docs: https://bit.ly/2VYPQD1
func (e *Email) sendEmail(ctx context.Context, req *pb.SendRequest) error {
	content := []interface{}{}
	replyTo := e.config.EmailFrom
	if len(req.ReplyTo) > 0 {
		replyTo = req.ReplyTo
	}

	if len(req.TextBody) > 0 {
		content = append(content, map[string]string{
			"type":  "text/plain",
			"value": req.TextBody,
		})
	}

	if len(req.HtmlBody) > 0 {
		content = append(content, map[string]string{
			"type":  "text/html",
			"value": req.HtmlBody,
		})
	}

	reqMap := map[string]interface{}{
		"from": map[string]string{
			"email": e.config.EmailFrom,
			"name":  req.From,
		},
		"reply_to": map[string]string{
			"email": replyTo,
		},
		"subject": req.Subject,
		"content": content,
		"personalizations": []interface{}{
			map[string]interface{}{
				"to": []map[string]string{
					{
						"email": req.To,
					},
				},
			},
		},
	}
	if len(e.config.PoolName) > 0 {
		reqMap["ip_pool_name"] = e.config.PoolName
	}

	reqBody, _ := json.Marshal(reqMap)

	httpReq, err := http.NewRequest("POST", "https://api.sendgrid.com/v3/mail/send", bytes.NewBuffer(reqBody))
	if err != nil {
		return err
	}

	httpReq.Header.Set("Authorization", "Bearer "+e.config.Key)
	httpReq.Header.Set("Content-Type", "application/json")

	rsp, err := new(http.Client).Do(httpReq)
	if err != nil {
		return fmt.Errorf("could not send email, error: %v", err)
	}
	defer rsp.Body.Close()

	tnt, ok := tenant.FromContext(ctx)
	if ok {
		msgID := rsp.Header.Get("X-Message-ID")
		if len(msgID) > 0 {
			sent := Sent{
				UserID:        tnt,
				SendgridMsgID: msgID,
			}
			b, _ := json.Marshal(&sent)
			if err := store.Write(&store.Record{
				Key:   fmt.Sprintf("%s/%s/%s", prefixUserID, sent.UserID, sent.SendgridMsgID),
				Value: b,
			}); err != nil {
				log.Errorf("Failed to persist mapping %+v %s", sent, err)
			}
			if err := store.Write(&store.Record{
				Key:   fmt.Sprintf("%s/%s", prefixSendgridID, sent.SendgridMsgID),
				Value: b,
			}); err != nil {
				log.Errorf("Failed to persist mapping %+v %s", sent, err)
			}
		}
	}

	if rsp.StatusCode < 200 || rsp.StatusCode > 299 {
		bytes, err := ioutil.ReadAll(rsp.Body)
		if err != nil {
			return err
		}
		return fmt.Errorf("could not send email, error: %v", string(bytes))
	}

	return nil
}

func (e *Email) Parse(ctx context.Context, req *pb.ParseRequest, rsp *pb.ParseResponse) error {
	if len(req.Address) == 0 {
		return errors.BadRequest("email.validate", "address is blank")
	}

	a, err := mail.ParseAddress(req.Address)
	if err != nil {
		return errors.InternalServerError("email.parse", err.Error())
	}

	rsp.Name = a.Name
	rsp.Address = a.Address

	return nil
}

func (e *Email) Validate(ctx context.Context, req *pb.ValidateRequest, rsp *pb.ValidateResponse) error {
	if len(req.Address) == 0 {
		return errors.BadRequest("email.validate", "address is blank")
	}

	_, err := mail.ParseAddress(req.Address)
	if err != nil {
		return nil
	}

	// success
	rsp.IsValid = true

	return nil
}
