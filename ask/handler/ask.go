package handler

import (
	"context"

	"github.com/ajanicij/goduckgo/goduckgo"
	"github.com/micro/micro/v3/service/errors"
	pb "github.com/micro/services/ask/proto"
)

type Ask struct{}

func (a *Ask) Question(ctx context.Context, req *pb.QuestionRequest, rsp *pb.QuestionResponse) error {
	if len(req.Query) == 0 {
		return errors.BadRequest("ask.question", "need a question")
	}

	msg, err := goduckgo.Query(req.Query)
	if err != nil {
		return errors.InternalServerError("ask.question", err.Error())
	}

	if len(msg.Abstract) > 0 {
		rsp.Answer = msg.Abstract
	} else if len(msg.AbstractText) > 0 {
		rsp.Answer = msg.AbstractText
	} else if len(msg.RelatedTopics) > 0 {
		rsp.Answer = msg.RelatedTopics[0].Text
	} else {
		rsp.Answer = "Sorry I don't know 😞"
		return nil
	}

	if len(msg.AbstractURL) > 0 {
		rsp.Url = msg.AbstractURL
	} else if len(msg.RelatedTopics) > 0 {
		rsp.Url = msg.RelatedTopics[0].FirstURL
	}

	if len(msg.Image) > 0 {
		rsp.Image = "https://duckduckgo.com" + msg.Image
	} else if len(msg.RelatedTopics) > 0 {
		rsp.Image = "https://duckduckgo.com" + msg.RelatedTopics[0].Icon.URL
	}

	return nil
}
