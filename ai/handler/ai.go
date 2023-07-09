package handler

import (
	"context"
	"time"

	pb "github.com/micro/services/ai/proto"
	"github.com/micro/services/pkg/api"
	"github.com/micro/services/pkg/tenant"
	"micro.dev/v4/service/config"
	"micro.dev/v4/service/errors"
	log "micro.dev/v4/service/logger"
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

func (e *Ai) Complete(ctx context.Context, req *pb.CompleteRequest, rsp *pb.CompleteResponse) error {
	if len(req.Text) == 0 {
		return errors.BadRequest("ai.complete", "missing text")
	}

	// get the tenant
	tnt, ok := tenant.FromContext(ctx)
	if !ok {
		tnt = "micro"
	}

	uri := "https://api.openai.com/v1/completions"

	var resp map[string]interface{}
	if err := api.Post(uri, map[string]interface{}{
		"model":       "text-davinci-003",
		"prompt":      req.Text,
		"max_tokens":  1000,
		"temperature": 0,
		"user":        tnt,
	}, &resp); err != nil {
		log.Errorf("Failed AI call: %v\n", err)
		return errors.InternalServerError("ai.complete", "Failed to make request")
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

func (e *Ai) Edit(ctx context.Context, req *pb.EditRequest, rsp *pb.EditResponse) error {
	if len(req.Text) == 0 {
		return errors.BadRequest("ai.edit", "missing text")
	}

	uri := "https://api.openai.com/v1/edits"

	if len(req.Instruction) == 0 {
		req.Instruction = "Edit the spelling and grammar"
	}

	var resp map[string]interface{}
	if err := api.Post(uri, map[string]interface{}{
		"model":       "text-davinci-edit-001",
		"input":       req.Text,
		"instruction": req.Instruction,
	}, &resp); err != nil {
		log.Errorf("Failed AI call: %v\n", err)
		return errors.InternalServerError("ai.edit", "Failed to make request")
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

func (e *Ai) Generate(ctx context.Context, req *pb.GenerateRequest, rsp *pb.GenerateResponse) error {
	if len(req.Text) == 0 {
		return errors.BadRequest("ai.generate", "missing image text")
	}

	// get the tenant
	tnt, ok := tenant.FromContext(ctx)
	if !ok {
		tnt = "micro"
	}

	uri := "https://api.openai.com/v1/images/generations"

	if req.Limit == 0 || req.Limit > 10 {
		req.Limit = 1
	}

	switch req.Size {
	case "256x256", "512x512", "1024x1024":
	default:
		req.Size = "1024x1024"
	}

	var resp map[string]interface{}
	if err := api.Post(uri, map[string]interface{}{
		"prompt":          req.Text,
		"n":               req.Limit,
		"size":            req.Size,
		"user":            tnt,
		"response_format": "b64_json",
	}, &resp); err != nil {
		log.Errorf("Failed AI Generate generation: %v\n", err)
		return errors.InternalServerError("ai.generate", "Failed to make request")
	}

	v := resp["data"]
	if v == nil {
		return nil
	}

	for _, i := range v.([]interface{}) {
		d := i.(map[string]interface{})
		rsp.Images = append(rsp.Images, &pb.Image{
			// TODO: upload image
			//Url: d["url"].(string),
			Base64: d["b64_json"].(string),
		})
	}

	return nil
}
