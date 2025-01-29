package handler

import (
	"context"
	"strings"

	"github.com/micro/micro/v5/service/errors"
	"github.com/micro/micro/v5/service/logger"
	pb "github.com/micro/services/google/proto"
	"google.golang.org/api/customsearch/v1"
	"google.golang.org/api/option"
)

type Google struct {
	Client *customsearch.Service
	CxId   string
}

func New(apiKey, cxId string) *Google {
	ctx := context.TODO()
	cs, _ := customsearch.NewService(ctx, option.WithAPIKey(apiKey))

	return &Google{
		Client: cs,
		CxId:   cxId,
	}
}

func (g *Google) Search(ctx context.Context, req *pb.SearchRequest, rsp *pb.SearchResponse) error {
	if len(req.Query) == 0 {
		return errors.BadRequest("google.search", "missing query")
	}

	resp, err := g.Client.Cse.List().Cx(g.CxId).Q(req.Query).Num(10).Do()
	if err != nil {
		logger.Errorf("failed to search google for %v: %v", req.Query, err)
		return errors.InternalServerError("google.search", "Failed to search for "+req.Query)
	}

	for _, item := range resp.Items {
		kind := strings.Split(item.Kind, "#")[1]
		rsp.Results = append(rsp.Results, &pb.SearchResult{
			Id:         item.CacheId,
			Kind:       kind,
			Title:      item.Title,
			Url:        item.Link,
			DisplayUrl: item.DisplayLink,
			Snippet:    item.Snippet,
		})
	}

	return nil
}
