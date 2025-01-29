package handler

import (
	"context"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/micro/micro/v5/service/config"
	"github.com/micro/micro/v5/service/errors"
	"github.com/micro/micro/v5/service/logger"
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

func (t *Twitter) Search(ctx context.Context, req *pb.SearchRequest, rsp *pb.SearchResponse) error {
	if len(req.Query) == 0 {
		return errors.BadRequest("twitter.query", "missing query")
	}

	if req.Limit <= 0 {
		req.Limit = 20
	}

	searchRsp, _, err := t.Client.Search.Tweets(&twitter.SearchTweetParams{
		Query: req.Query,
		Count: int(req.Limit),
	})
	if err != nil {
		logger.Errorf("Failed to retrieve tweets for %v: %v", req.Query, err)
		return errors.InternalServerError("twitter.search", "Failed to retrieve tweets for %v: %v", req.Query, err)
	}

	for _, tweet := range searchRsp.Statuses {
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

func (t *Twitter) User(ctx context.Context, req *pb.UserRequest, rsp *pb.UserResponse) error {
	if len(req.Username) == 0 {
		return errors.BadRequest("twitter.user", "missing username")
	}

	user, _, err := t.Client.Users.Show(&twitter.UserShowParams{
		ScreenName: req.Username,
	})
	if err != nil {
		logger.Errorf("Failed to retrieve user profile for %v: %v", req.Username, err)
		return errors.InternalServerError("twitter.user", "Failed to retrieve user profile for %v: %v", req.Username, err)
	}

	rsp.Status = &pb.Tweet{
		Id:              user.Status.ID,
		Text:            user.Status.Text,
		CreatedAt:       user.Status.CreatedAt,
		FavouritedCount: int64(user.Status.FavoriteCount),
		RetweetedCount:  int64(user.Status.RetweetCount),
		Username:        req.Username,
	}

	rsp.Profile = &pb.Profile{
		Id:          user.ID,
		Name:        user.Name,
		Username:    user.ScreenName,
		CreatedAt:   user.CreatedAt,
		Description: user.Description,
		Followers:   int64(user.FollowersCount),
		Private:     user.Protected,
		Verified:    user.Verified,
		ImageUrl:    user.ProfileImageURL,
	}

	return nil
}

func (t *Twitter) Trends(ctx context.Context, req *pb.TrendsRequest, rsp *pb.TrendsResponse) error {
	trendRsp, _, err := t.Client.Trends.Place(1, &twitter.TrendsPlaceParams{WOEID: 1})
	if err != nil {
		logger.Errorf("Failed to retrieve trends: %v", err)
		return errors.InternalServerError("twitter.trends", "Failed to retrieve trends")
	}

	for _, list := range trendRsp {
		for _, trend := range list.Trends {
			rsp.Trends = append(rsp.Trends, &pb.Trend{
				Name:        trend.Name,
				Url:         trend.URL,
				TweetVolume: trend.TweetVolume,
			})
		}
	}

	return nil
}
