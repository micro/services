package main

import (
	posts "github.com/micro/services/posts/proto"
	"github.com/micro/services/sentiment/handler"
	pb "github.com/micro/services/sentiment/proto"
	"github.com/micro/services/sentiment/subscriber"

	"github.com/micro/micro/v3/service"
	"github.com/micro/micro/v3/service/logger"
)

func main() {
	// Create service
	srv := service.New(
		service.Name("sentiment"),
	)

	// Register handler
	pb.RegisterSentimentHandler(srv.Server(), new(handler.Sentiment))

	// Register subscriber
	service.Subscribe("posts", subscriber.EnrichPost)
	subscriber.PostsClient = posts.NewPostsService("posts", srv.Client())

	// Run service
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
