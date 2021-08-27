package handler

import (
	"context"
	"net/url"

	"github.com/enescakir/emoji"
	"github.com/kevinburke/twilio-go"
	"github.com/micro/micro/v3/service/config"
	"github.com/micro/micro/v3/service/errors"
	"github.com/micro/micro/v3/service/logger"
	pb "github.com/micro/services/emoji/proto"
)

type Emoji struct{}

func (e *Emoji) Find(ctx context.Context, req *pb.FindRequest, rsp *pb.FindResponse) error {
	emoji, ok := emoji.Find(req.Alias)
	if !ok {
		return errors.NotFound("emoji.find", req.Alias+" not found")
	}
	rsp.Emoji = emoji
	return nil
}

func (e *Emoji) Flag(ctx context.Context, req *pb.FlagRequest, rsp *pb.FlagResponse) error {
	emoji, err := emoji.CountryFlag(req.Code)
	if err != nil {
		return errors.BadRequest("emoji.flag", err.Error())
	}
	rsp.Flag = emoji.String()
	return nil
}

func (e *Emoji) Print(ctx context.Context, req *pb.PrintRequest, rsp *pb.PrintResponse) error {
	rsp.Text = emoji.Parse(req.Text)
	return nil
}

func (e *Emoji) Send(ctx context.Context, req *pb.SendRequest, rsp *pb.SendResponse) error {
	if len(req.From) == 0 {
		return errors.BadRequest("emoji.send", "require from field")
	}
	if len(req.To) == 0 {
		return errors.BadRequest("emoji.send", "require to field")
	}
	if len(req.Message) == 0 {
		return errors.BadRequest("emoji.send", "message is blank")
	}

	v, err := config.Get("twilio.sid")
	if err != nil {
		logger.Error("Failed to get twilio.sid config")
		return errors.InternalServerError("emoji.send", "failed to send message")
	}
	sid := v.String("")

	v, err = config.Get("twilio.token")
	if err != nil {
		logger.Error("Failed to get twilio.token config")
		return errors.InternalServerError("emoji.send", "failed to send message")
	}
	token := v.String("")

	v, err = config.Get("twilio.number")
	if err != nil {
		logger.Error("Failed to get twilio.number config")
		return errors.InternalServerError("emoji.send", "failed to send message")
	}
	number := v.String("")

	message := emoji.Parse(req.Message)
	message += "  Sent from " + req.From

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
		return errors.InternalServerError("emoji.send", "failed to send message: %v", err.Error())
	}

	rsp.Success = true

	return nil
}
