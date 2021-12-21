package handler

import (
	"context"

	"github.com/micro/micro/v3/service/config"
	me "github.com/micro/micro/v3/service/errors"
	"github.com/micro/micro/v3/service/logger"
	"github.com/pkg/errors"
	"golang.org/x/text/language"
	"google.golang.org/api/option"

	pb "github.com/micro/services/translate/proto"

	"cloud.google.com/go/translate"
)

type translation struct {
	ApiKey string
	Limit  int
}

func NewTranslation() *translation {
	v, err := config.Get("translate.google.api_key")
	if err != nil {
		logger.Fatalf("translate.google.api_key config not found: %v", err)
	}
	key := v.String("")

	if key == "" {
		logger.Fatalf("translate.google.api_key config can not be an empty string")
	}

	v, err = config.Get("translate.text.char_limit")
	if err != nil {
		logger.Fatalf("translate.text.char_limit config not found: %v", err)
	}
	limit := v.Int(0)

	return &translation{
		ApiKey: key,
		Limit:  limit,
	}
}

// Text calls Google Cloud Translation Basic edition API
// For more information: https://cloud.google.com/translate/docs/samples/translate-text-with-model
func (t *translation) Text(ctx context.Context, req *pb.TextRequest, rsp *pb.TextResponse) error {
	client, err := translate.NewClient(ctx, option.WithAPIKey(t.ApiKey))
	if err != nil {
		return errors.Wrap(err, "new google translation client error")
	}
	defer client.Close()

	source, err := language.Parse(req.Source)
	if err != nil {
		return errors.Wrap(err, "google translation parse source language error")
	}

	target, err := language.Parse(req.Target)
	if err != nil {
		return errors.Wrap(err, "google translation parse target language error")
	}

	// TODO: configurable char limit
	if t.Limit > 0 && len(req.Content) > t.Limit {
		return me.BadRequest("google.translate", "Exceeds char limit %d", t.Limit)
	}

	result, err := client.Translate(ctx, []string{req.Content}, target, &translate.Options{
		Source: source,
		Format: translate.Format(req.Format),
		Model:  req.Model,
	})

	if err != nil {
		return errors.Wrap(err, "get google translation error")
	}

	if len(result) == 0 {
		return nil
	}

	rsp.Translation = &pb.Translation{
		Text:   result[0].Text,
		Source: result[0].Source.String(),
		Model:  result[0].Model,
	}

	return nil
}
