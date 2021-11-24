package handler

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/micro/micro/v3/service/errors"
	pb "github.com/micro/services/news/proto"
	"github.com/micro/services/pkg/api"
)

type News struct {
	apiKey string
}

type Article struct {
	UUID        string     `json:"uuid,omitempty"`
	Title       string     `json:"title,omitempty"`
	Description string     `json:"description,omitempty"`
	Keywords    string     `json:"keywords,omitempty"`
	Snippet     string     `json:"snippet,omitempty"`
	Url         string     `json:"url,omitempty"`
	ImageUrl    string     `json:"image_url,omitempty"`
	PublishedAt string     `json:"published_at,omitempty"`
	Source      string     `json:"source,omitempty"`
	Categories  []string   `json:"categories,omitempty"`
	Language    string     `json:"language,omitempty"`
	Locale      string     `json:"locale,omitempty"`
	Similar     []*Article `json:"similar,omitempty"`
}

type Headlines struct {
	Data map[string][]*Article `json:"data"`
}

var (
	apiURL = "https://api.thenewsapi.com"
)

func New(apiKey string) *News {
	return &News{
		apiKey: apiKey,
	}
}

func toProto(v *Article) *pb.Article {
	return &pb.Article{
		Id:          v.UUID,
		Title:       v.Title,
		Description: v.Description,
		Keywords:    v.Keywords,
		Snippet:     v.Snippet,
		Url:         v.Url,
		ImageUrl:    v.ImageUrl,
		PublishedAt: v.PublishedAt,
		Language:    v.Language,
		Source:      v.Source,
		Categories:  v.Categories,
		Locale:      v.Locale,
	}
}

func (n *News) Headlines(ctx context.Context, req *pb.HeadlinesRequest, rsp *pb.HeadlinesResponse) error {
	path := "/v1/news/headlines"
	locale := "us"
	language := "en"
	date := time.Now().Format("2006-01-02")

	if len(req.Locale) > 0 {
		locale = req.Locale
	}
	if len(req.Language) > 0 {
		language = req.Language
	}

	if len(req.Date) > 0 {
		date = req.Date
	}

	vals := url.Values{}
	vals.Set("api_token", n.apiKey)
	vals.Set("locale", locale)
	vals.Set("published_on", date)
	vals.Set("language", language)

	uri := fmt.Sprintf("%s%s?%s", apiURL, path, vals.Encode())
	var resp *Headlines
	if err := api.Get(uri, &resp); err != nil {
		return errors.InternalServerError("news.headlines", err.Error())
	}

	for _, articles := range resp.Data {
		for _, v := range articles {
			rsp.Articles = append(rsp.Articles, toProto(v))

			for _, a := range v.Similar {
				rsp.Articles = append(rsp.Articles, toProto(a))
			}
		}
	}

	return nil
}
