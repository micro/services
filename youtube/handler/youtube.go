package handler

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/micro/micro/v3/service/errors"
	"github.com/micro/micro/v3/service/logger"
	pb "github.com/micro/services/youtube/proto"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

type Youtube struct {
	Client *youtube.Service
}

var (
	iframe = `<iframe width="560" height="315" src="%s" title="YouTube video player" frameborder="0" allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture" allowfullscreen></iframe>`
)

func New(apiKey string) *Youtube {
	ctx := context.TODO()
	yt, _ := youtube.NewService(ctx, option.WithAPIKey(apiKey))

	return &Youtube{
		Client: yt,
	}
}

func (y *Youtube) Embed(ctx context.Context, req *pb.EmbedRequest, rsp *pb.EmbedResponse) error {
	if len(req.Url) == 0 {
		return errors.BadRequest("youtube.embed", "missing url")
	}

	if strings.HasPrefix(req.Url, "https://youtu.be/") {
		id := strings.TrimPrefix(req.Url, "https://youtu.be/")

		rsp.Link = "https://www.youtube.com/embed/" + id
		rsp.Script = fmt.Sprintf(iframe, rsp.Link)
		rsp.ShortUrl = req.Url

		return nil
	}

	if !strings.HasPrefix(req.Url, "https://www.youtube.com/watch") {
		return errors.BadRequest("youtube.embed", "invalid url")
	}

	uri, err := url.Parse(req.Url)
	if err != nil {
		return errors.BadRequest("youtube.embed", "invalid url")
	}

	vals := uri.Query()
	id := vals.Get("v")

	if len(id) == 0 {
		return errors.BadRequest("youtube.embed", "invalid url")
	}

	rsp.Link = "https://www.youtube.com/embed/" + id
	rsp.Script = fmt.Sprintf(iframe, rsp.Link)
	rsp.ShortUrl = "https://youtu.be/" + id

	return nil
}

func (y *Youtube) Search(ctx context.Context, req *pb.SearchRequest, rsp *pb.SearchResponse) error {
	if len(req.Query) == 0 {
		return errors.BadRequest("youtube.search", "missing query")
	}

	resp, err := y.Client.Search.List([]string{"id", "snippet"}).Q(req.Query).MaxResults(25).Do()
	if err != nil {
		logger.Errorf("failed to search youtube for %v: %v", req.Query, err)
		return errors.InternalServerError("youtube.search", "Failed to search for "+req.Query)
	}

	for _, item := range resp.Items {
		var id, url string
		kind := strings.Split(item.Id.Kind, "#")[1]
		switch kind {
		case "video":
			id = item.Id.VideoId
			url = "https://www.youtube.com/watch?v=" + id
		case "playlist":
			id = item.Id.PlaylistId
			url = "https://www.youtube.com/playlist?list=" + id
		case "channel":
			id = item.Id.ChannelId
			url = "https://www.youtube.com/channel/" + id
		}
		rsp.Results = append(rsp.Results, &pb.SearchResult{
			Id:           id,
			Kind:         kind,
			Title:        item.Snippet.Title,
			ChannelId:    item.Snippet.ChannelId,
			ChannelTitle: item.Snippet.ChannelTitle,
			Description:  item.Snippet.Description,
			PublishedAt:  item.Snippet.PublishedAt,
			Broadcasting: item.Snippet.LiveBroadcastContent,
			Url:          url,
		})
	}

	return nil
}
