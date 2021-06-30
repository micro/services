package handler

import (
	"context"
	"net/url"

	"github.com/kevinburke/twilio-go"
	"github.com/micro/micro/v3/service/config"
	"github.com/micro/micro/v3/service/errors"
	"github.com/micro/micro/v3/service/logger"
	pb "github.com/micro/services/sms/proto"
)

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

	vals := url.Values{}
	vals.Set("Body", message)
	vals.Set("From", number)
	vals.Set("To", req.To)
	// non configurable and must match publicapi.json
	vals.Set("MaxPrice", "0.01")

	client := twilio.NewClient(sid, token, nil)
	_, err = client.Messages.Create(ctx, vals)
	if err != nil {
		logger.Errorf("Failed to send message: %v", err)
		return errors.InternalServerError("sms.send", "failed to send message: %v", err.Error())
	}

	rsp.Status = "ok"

	return nil
}
