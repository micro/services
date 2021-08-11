package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/kevinburke/twilio-go"
	"github.com/micro/micro/v3/service/auth"
	"github.com/micro/micro/v3/service/config"
	"github.com/micro/micro/v3/service/errors"
	"github.com/micro/micro/v3/service/logger"
	"github.com/micro/micro/v3/service/store"
	"github.com/micro/services/pkg/tenant"
	pb "github.com/micro/services/sms/proto"
)

const (
	prefixUserID   = "byUserID"
	prefixTwilioID = "byTwilioID"
)

type Sent struct {
	UserID      string
	TwilioMsgID string
}

type Sms struct{}

func (e *Sms) Send(ctx context.Context, req *pb.SendRequest, rsp *pb.SendResponse) error {
	if len(req.From) == 0 {
		return errors.BadRequest("sms.send", "require from field")
	}
	if len(req.To) == 0 {
		return errors.BadRequest("sms.send", "require to field")
	}
	if len(req.Message) == 0 {
		return errors.BadRequest("sms.send", "message is blank")
	}

	tnt, _ := tenant.FromContext(ctx)
	// crudely ban any sender in the banned list aka no impersonating
	frm := strings.ToLower(req.From)
	for _, sender := range BanFrom {
		if strings.Contains(frm, strings.ToLower(sender)) {
			acc, _ := auth.AccountFromContext(ctx)

			logger.Error("Request to send from %v blocked by account: %v tenant: %v", req.From, acc, tnt)
			return errors.BadRequest("sms.send", "sender blocked")
		}
	}

	v, err := config.Get("twilio.sid")
	if err != nil {
		logger.Error("Failed to get twilio.sid config")
		return errors.InternalServerError("sms.send", "failed to send message")
	}
	sid := v.String("")

	v, err = config.Get("twilio.token")
	if err != nil {
		logger.Error("Failed to get twilio.token config")
		return errors.InternalServerError("sms.send", "failed to send message")
	}
	token := v.String("")

	v, err = config.Get("twilio.number")
	if err != nil {
		logger.Error("Failed to get twilio.number config")
		return errors.InternalServerError("sms.send", "failed to send message")
	}
	number := v.String("")

	message := req.Message + "  Sent from " + req.From

	if len(message) > 160 {
		return errors.BadRequest("sms.send", "message is too long")
	}

	vals := url.Values{}
	vals.Set("Body", message)
	vals.Set("From", number)
	vals.Set("To", req.To)
	// non configurable and must match publicapi.json
	vals.Set("MaxPrice", "0.01")

	client := twilio.NewClient(sid, token, nil)
	twMsg, err := client.Messages.Create(ctx, vals)
	if err != nil {
		logger.Errorf("Failed to send message: %v", err)
		return errors.InternalServerError("sms.send", "failed to send message: %v", err.Error())
	}

	sent := Sent{
		UserID:      tnt,
		TwilioMsgID: twMsg.Sid,
	}

	b, _ := json.Marshal(&sent)
	// log the association so we can correlate who sent what
	if err := store.Write(&store.Record{
		Key:   fmt.Sprintf("%s/%s/%s/%s", prefixUserID, sent.UserID, time.Now().Format("20060102"), sent.TwilioMsgID),
		Value: b,
	}); err != nil {
		logger.Errorf("Failed to store association %+v %s", sent, err)
	}
	if err := store.Write(&store.Record{
		Key:   fmt.Sprintf("%s/%s", prefixTwilioID, sent.TwilioMsgID),
		Value: b,
	}); err != nil {
		logger.Errorf("Failed to store association %+v %s", sent, err)
	}

	rsp.Status = "ok"

	return nil
}
