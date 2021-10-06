package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/micro/micro/v3/service/config"
	"github.com/micro/micro/v3/service/errors"
	"github.com/micro/micro/v3/service/logger"
	gifs "github.com/micro/services/gifs/proto"
)

const (
	giphySearch         = "https://api.giphy.com/v1/gifs/search?api_key=%s&q=%s&limit=%d&offset=%d&rating=%s&lang=%s"
	defaultLimit  int32 = 25
	defaultOffset int32 = 0
	defaultRating       = "g"
	defaultLang         = "en"
)

type conf struct {
	APIKey string `json:"api_key"`
}

type Gifs struct {
	conf conf
}

func New() *Gifs {
	v, err := config.Get("micro.gifs")
	if err != nil {
		logger.Fatalf("Failed to load config %s", err)
	}
	var c conf
	if err := v.Scan(&c); err != nil {
		logger.Fatalf("Failed to load config %s", err)
	}

	return &Gifs{conf: c}
}

func (g *Gifs) Search(ctx context.Context, request *gifs.SearchRequest, response *gifs.SearchResponse) error {
	if len(request.Query) == 0 {
		return errors.BadRequest("gifs.Search", "Missing query field")
	}
	limit := defaultLimit
	if request.Limit > 0 {
		limit = request.Limit
	}
	offset := defaultOffset
	if request.Offset > 0 {
		offset = request.Offset
	}

	rating := defaultRating
	if len(request.Rating) > 0 {
		rating = request.Rating
	}
	lan := defaultLang
	if len(request.Lang) > 0 {
		lan = request.Lang
	}
	rsp, err := http.Get(fmt.Sprintf(giphySearch, g.conf.APIKey, request.Query, limit, offset, rating, lan))
	if err != nil {
		logger.Errorf("Error querying giphy %s", err)
		return errors.InternalServerError("gifs.Search", "Error querying for gifs")
	}
	defer rsp.Body.Close()
	b, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		logger.Errorf("Error marshalling giphy response %s", err)
		return errors.InternalServerError("gifs.Search", "Error querying for gifs")
	}
	var gRsp searchResponse
	if err := json.Unmarshal(b, &gRsp); err != nil {
		logger.Errorf("Error marshalling giphy response %s", err)
		return errors.InternalServerError("gifs.Search", "Error querying for gifs")
	}
	response.Data = marshalGifs(gRsp.Data)
	response.Pagination = &gifs.Pagination{
		Offset:     gRsp.Pagination.Offset,
		TotalCount: gRsp.Pagination.TotalCount,
		Count:      gRsp.Pagination.Count,
	}
	return nil
}

func marshalGifs(in []gif) []*gifs.Gif {
	ret := make([]*gifs.Gif, len(in))
	for i, v := range in {
		ret[i] = &gifs.Gif{
			Id:       v.ID,
			Slug:     v.Slug,
			Url:      v.URL,
			ShortUrl: v.ShortURL,
			EmbedUrl: v.EmbedURL,
			Source:   v.Source,
			Rating:   v.Rating,
			Title:    v.Title,
			Images: &gifs.ImageFormats{
				Original:               marshalImageFormat(v.Images.Original),
				Downsized:              marshalImageFormat(v.Images.Downsized),
				FixedHeight:            marshalImageFormat(v.Images.FixedHeight),
				FixedHeightStill:       marshalImageFormat(v.Images.FixedHeightStill),
				FixedHeightDownsampled: marshalImageFormat(v.Images.FixedHeightDownsampled),
				FixedWidth:             marshalImageFormat(v.Images.FixedWidth),
				FixedWidthStill:        marshalImageFormat(v.Images.FixedWidthStill),
				FixedWidthDownsampled:  marshalImageFormat(v.Images.FixedWidthDownsampled),
				FixedHeightSmall:       marshalImageFormat(v.Images.FixedHeightSmall),
				FixedHeightSmallStill:  marshalImageFormat(v.Images.FixedHeightSmallStill),
				FixedWidthSmall:        marshalImageFormat(v.Images.FixedWidthSmall),
				FixedWidthSmallStill:   marshalImageFormat(v.Images.FixedWidthSmallStill),
				DownsizedStill:         marshalImageFormat(v.Images.DownsizedStill),
				DownsizedLarge:         marshalImageFormat(v.Images.DownsizedLarge),
				DownsizedMedium:        marshalImageFormat(v.Images.DownsizedMedium),
				DownsizedSmall:         marshalImageFormat(v.Images.DownsizedSmall),
				OriginalStill:          marshalImageFormat(v.Images.OriginalStill),
				Looping:                marshalImageFormat(v.Images.Looping),
				Preview:                marshalImageFormat(v.Images.Preview),
				PreviewGif:             marshalImageFormat(v.Images.PreviewGif),
			},
		}

	}
	return ret
}

func marshalImageFormat(in format) *gifs.ImageFormat {
	mustInt32 := func(s string) int32 {
		i, _ := strconv.Atoi(s)
		return int32(i)
	}
	return &gifs.ImageFormat{
		Height:   mustInt32(in.Height),
		Width:    mustInt32(in.Width),
		Size:     mustInt32(in.Size),
		Url:      in.URL,
		Mp4Url:   in.MP4URL,
		Mp4Size:  mustInt32(in.MP4Size),
		WebpUrl:  in.WebpURL,
		WebpSize: mustInt32(in.WebpSize),
	}
}
