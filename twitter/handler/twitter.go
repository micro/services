package handler

import (
	"context"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/micro/micro/v3/service/config"
	"github.com/micro/micro/v3/service/errors"
	"github.com/micro/micro/v3/service/logger"
	pb "github.com/micro/services/twitter/proto"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

type Twitter struct {
	Client *twitter.Client
}

func New() *Twitter {
	v, err := config.Get("twitter.api_key")
	if err != nil {
		logger.Fatalf("twitter.api_key config not found: %v", err)
	}
	apiKey := v.String("")
	if len(apiKey) == 0 {
		logger.Fatal("twitter.api_key config not found")
	}
	v, err = config.Get("twitter.api_secret")
	if err != nil {
		logger.Fatalf("twitter.api_secret config not found: %v", err)
	}
	apiSecret := v.String("")
	if len(apiSecret) == 0 {
		logger.Fatal("twitter.api_secret config not found")
	}

	// oauth2 configures a client that uses app credentials to keep a fresh token
	config := &clientcredentials.Config{
		ClientID:     apiKey,
		ClientSecret: apiSecret,
		TokenURL:     "https://api.twitter.com/oauth2/token",
	}
	// http.Client will automatically authorize Requests
	httpClient := config.Client(oauth2.NoContext)

	// Twitter client
	client := twitter.NewClient(httpClient)

	return &Twitter{
		Client: client,
	}
}

func (t *Twitter) Timeline(ctx context.Context, req *pb.TimelineRequest, rsp *pb.TimelineResponse) error {
	if len(req.Username) == 0 {
		return errors.BadRequest("twitter.timeline", "missing username")
	}

	if req.Limit <= 0 {
		req.Limit = 20
	}

	tweets, _, err := t.Client.Timelines.UserTimeline(&twitter.UserTimelineParams{
		ScreenName: req.Username,
		Count:      int(req.Limit),
	})
	if err != nil {
		logger.Errorf("Failed to retrieve tweets for %v: %v", req.Username, err)
		return errors.InternalServerError("twitter.timeline", "Failed to retrieve tweets for %v: %v", req.Username, err)
	}

	for _, tweet := range tweets {
		rsp.Tweets = append(rsp.Tweets, &pb.Tweet{
			Id:              tweet.ID,
			Text:            tweet.Text,
			CreatedAt:       tweet.CreatedAt,
			FavouritedCount: int64(tweet.FavoriteCount),
			RetweetedCount:  int64(tweet.RetweetCount),
			Username:        tweet.User.ScreenName,
		})
	}

	return nil
}
