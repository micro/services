package handler

import (
	"context"
	"fmt"
	"net/url"

	"github.com/micro/micro/v3/service/errors"
	pb "github.com/micro/services/twitter/proto"
	"github.com/micro/services/twitter/api"
)

type Api struct{}

func (a *Api) Tweet(ctx context.Context, req *pb.TweetRequest, rsp *pb.TweetResponse) error {
	if len(req.Status) == 0 {
		return errors.BadRequest("twitter.Api.Tweet", "Status cannot be blank")
	}
	if len(req.Status) > 140 {
		return errors.BadRequest("twitter.Api.Tweet", "Status cannot be longer than 140 chars")
	}

	u := url.Values{}

	if req.InReplyToStatusId > 0 {
		u.Set("in_reply_to_status_id", fmt.Sprintf("%d", req.InReplyToStatusId))
	}

	if req.PossiblySensitive {
		u.Set("possibly_sensitive", "true")
	}

	if req.LatLng != nil {
		u.Set("lat", fmt.Sprintf("%.6f", req.LatLng.Lat))
		u.Set("long", fmt.Sprintf("%.6f", req.LatLng.Lng))
	}

	if len(req.PlaceId) > 0 {
		u.Set("place_id", req.PlaceId)
	}

	if req.DisplayCoordinates {
		u.Set("display_coordinates", "true")
	}

	if req.TrimUser {
		u.Set("trim_user", "true")
	}

	for _, id := range req.MediaIds {
		u.Set("media_ids", fmt.Sprintf("%d", id))
	}

	if err := api.Tweet(req.Status, u, rsp); err != nil {
		return errors.InternalServerError("twitter.Api.Tweet", err.Error())
	}

	return nil
}
