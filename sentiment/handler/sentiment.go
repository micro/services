package handler

import (
	"context"

	"github.com/cdipaolo/sentiment"
	"github.com/micro/micro/v3/service/errors"
	pb "github.com/micro/services/sentiment/proto"
)

type Sentiment struct {
	Model *sentiment.Models
}

func (e *Sentiment) Analyze(ctx context.Context, req *pb.Request, rsp *pb.Response) error {
	if len(req.Text) == 0 {
		return errors.BadRequest("sentiment.analyze", "text is blank")
	}

	if len(req.Lang) == 0 {
		req.Lang = "english"
	}

	if req.Lang != "english" {
		return errors.BadRequest("sentiment.analyze", "only support english")
	}

	an := e.Model.SentimentAnalysis(req.Text, sentiment.English)
	rsp.Score = float64(an.Score)

	// TODO: more complex word scoring

	return nil
}
