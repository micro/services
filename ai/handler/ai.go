package handler

import (
	"context"
	"time"

	"github.com/micro/micro/v3/service/config"
	"github.com/micro/micro/v3/service/errors"
	log "github.com/micro/micro/v3/service/logger"
	pb "github.com/micro/services/ai/proto"
	"github.com/micro/services/pkg/api"
	"github.com/micro/services/pkg/tenant"
)

type Ai struct{}

// Return a new handler
func New() *Ai {
	v, err := config.Get("ai.api_key")
	if err != nil {
		log.Fatal(err)
	}
	key := v.String("")
	if len(key) == 0 {
		log.Fatal("Missing api key")
	}
	api.SetKey("Authorization", "Bearer "+key)
	api.SetCache(true, time.Minute*10)

	return &Ai{}
}

func (e *Ai) Call(ctx context.Context, req *pb.Request, rsp *pb.Response) error {
	if len(req.Text) == 0 {
		return errors.BadRequest("ai.call", "missing text")
	}

	// get the tenant
	tnt, ok := tenant.FromContext(ctx)
	if !ok {
		tnt = "micro"
	}

	uri := "https://api.openai.com/v1/completions"

	var resp map[string]interface{}
	if err := api.Post(uri, map[string]interface{}{
		"model":       "text-davinci-002",
		"prompt":      req.Text,
		"max_tokens":  1000,
		"temperature": 0,
		"user":        tnt,
	}, &resp); err != nil {
		log.Errorf("Failed AI call: %v\n", err)
		return errors.InternalServerError("ai.call", "Failed to make request")
	}

	v := resp["choices"]
	if v == nil {
		return nil
	}

	// get first choice
	choice := v.([]interface{})[0].(map[string]interface{})

	// set response text
	rsp.Text = choice["text"].(string)

	return nil
}

func (e *Ai) Moderate(ctx context.Context, req *pb.ModerateRequest, rsp *pb.ModerateResponse) error {
	if len(req.Text) == 0 {
		return errors.BadRequest("ai.moderate", "missing text")
	}

	uri := "https://api.openai.com/v1/moderations"

	var resp map[string]interface{}
	if err := api.Post(uri, map[string]interface{}{
		"input": req.Text,
	}, &resp); err != nil {
		log.Errorf("Failed AI moderation: %v\n", err)
		return errors.InternalServerError("ai.moderate", "Failed to make request")
	}

	v := resp["results"]
	if v == nil {
		return nil
	}

	// get first choice
	results := v.([]interface{})[0].(map[string]interface{})

	// set response text
	rsp.Flagged, _ = results["flagged"].(bool)

	rsp.Categories = make(map[string]bool)
	rsp.Scores = make(map[string]float64)

	// set the categories
	for k, v := range results["categories"].(map[string]interface{}) {
		rsp.Categories[k] = v.(bool)
	}

	// set the scores
	for k, v := range results["category_scores"].(map[string]interface{}) {
		rsp.Scores[k] = v.(float64)
	}

	return nil
}
