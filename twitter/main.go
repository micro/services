package main

import (
	"github.com/micro/micro/v3/service/logger"
	"github.com/micro/micro/v3/service/config"
	"github.com/micro/micro/v3/service"
	"github.com/micro/services/twitter/handler"
	pb "github.com/micro/services/twitter/proto"
	"github.com/micro/services/twitter/api"
)

func getVal(v string) string {
	val, err := config.Get(v)
	if err != nil {
		return ""
	}
	return val.String("")
}

func configure() {
	api.Token = getVal("twitter.access_token")
	api.TokenSecret = getVal("twitter.access_token_secret")
	api.ConsumerKey = getVal("twitter.consumer_key")
	api.ConsumerSecret = getVal("twitter.consumer_secret")
	api.Init()
}

func init() {
	configure()
}

func main() {
	service := service.New(
		service.Name("twitter"),
	)

	service.Init()

	pb.RegisterApiHandler(service.Server(), &handler.Api{})

	if err := service.Run(); err != nil {
		logger.Fatal(err)
	}
}
