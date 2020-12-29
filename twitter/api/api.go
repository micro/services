// Package api is a simple twitter api client
package api

import (
	"encoding/json"
	"net/url"

	"github.com/ChimeraCoder/anaconda"
)

var (
	Token          string
	TokenSecret    string
	ConsumerKey    string
	ConsumerSecret string

	api *anaconda.TwitterApi
)

func Init() {
	anaconda.SetConsumerKey(ConsumerKey)
	anaconda.SetConsumerSecret(ConsumerSecret)
	api = anaconda.NewTwitterApi(Token, TokenSecret)
}

func Tweet(status string, args url.Values, rsp interface{}) error {
	tweet, err := api.PostTweet(status, args)
	if err != nil {
		return err
	}

	b, err := json.Marshal(tweet)
	if err != nil {
		return err
	}

	return json.Unmarshal(b, rsp)
}
