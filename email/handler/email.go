package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/micro/micro/v3/service/config"
	"github.com/micro/micro/v3/service/errors"
	log "github.com/micro/micro/v3/service/logger"
	"github.com/micro/micro/v3/service/store"
	pb "github.com/micro/services/email/proto"
	"github.com/micro/services/pkg/tenant"
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
}

func NewEmailHandler() *Email {
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
	}
}

type Email struct {
	config sendgridConf
}

func (e *Email) Send(ctx context.Context, request *pb.SendRequest, response *pb.SendResponse) error {
	if len(request.From) == 0 {
		return errors.BadRequest("email.send.validation", "Missing from address")
	}
	if len(request.To) == 0 {
		return errors.BadRequest("email.send.validation", "Missing to address")
	}
	if len(request.Subject) == 0 {
		return errors.BadRequest("email.send.validation", "Missing subject")
	}
	if len(request.TextBody) == 0 && len(request.HtmlBody) == 0 {
		return errors.BadRequest("email.send.validation", "Missing email body")
	}

	if err := e.sendEmail(ctx, request); err != nil {
		log.Errorf("Error sending email: %v\n", err)
		return errors.InternalServerError("email.sendemail", "Error sending email")
	}

	return nil
}

// sendEmail sends an email invite via the sendgrid API using the
// pre-designed email template. Docs: https://bit.ly/2VYPQD1
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

	reqBody, _ := json.Marshal(map[string]interface{}{
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
	})

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
