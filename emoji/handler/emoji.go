package handler

import (
	"context"

	"github.com/enescakir/emoji"
	"github.com/micro/micro/v3/service/errors"
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

