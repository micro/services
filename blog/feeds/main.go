package main

import (
	pb "github.com/micro/services/blog/feeds/proto"
	posts "github.com/micro/services/blog/posts/proto"

	"github.com/micro/services/blog/feeds/handler"

	"github.com/micro/micro/v3/service"
	"github.com/micro/micro/v3/service/logger"
)

func main() {
	// Create service
	srv := service.New(
		service.Name("feeds"),
		service.Version("latest"),
	)

	// Register handler
	pb.RegisterFeedsHandler(srv.Server(), handler.NewFeeds(posts.NewPostsService("posts", srv.Client())))

	// Run service
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
