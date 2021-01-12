package main

import (
	comments "github.com/micro/services/comments/proto"
	"github.com/micro/services/blog/handler"
	posts "github.com/micro/services/posts/proto"
	proto "github.com/micro/services/blog/proto"
	tags "github.com/micro/services/tags/proto"

	"github.com/micro/micro/v3/service"
	"github.com/micro/micro/v3/service/logger"
)

func main() {
	// Create service
	srv := service.New(
		service.Name("blog"),
	)

	// Register handler
	proto.RegisterBlogHandler(srv.Server(), handler.NewBlog(
		posts.NewPostsService("posts", srv.Client()),
		comments.NewCommentsService("comments", srv.Client()),
		tags.NewTagsService("tags", srv.Client()),
	))

	// Run service
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
