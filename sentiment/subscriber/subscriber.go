package subscriber

import (
	"context"
	"fmt"
	"time"

	pb "github.com/micro/services/posts/proto"
	"github.com/micro/services/sentiment/model"
)

var (
	// assume this is initialised by main
	PostsClient pb.PostsService
)

// EnrichPost will enrich a post with the sentiment and save it
func EnrichPost(ctx context.Context, post *pb.Post) error {
	if PostsClient == nil {
		return nil
	}

	// start by analysing the title
	// later we will look at the content
	score := model.Analyze(post.Title)
	post.Metadata["sentiment"] = fmt.Sprintf("%.1f", score)

	// now save the post
	PostsClient.Save(ctx, &pb.SaveRequest{
		Id:        post.Id,
		Title:     post.Title,
		Content:   post.Content,
		Timestamp: time.Now().Unix(),
		Metadata:  post.Metadata,
		Tags:      post.Tags,
		Image:     post.Image,
		Slug:      post.Slug,
	})

	return nil
}
