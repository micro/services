package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/micro/micro/v3/service/errors"
	"github.com/micro/micro/v3/service/logger"
	"github.com/micro/micro/v3/service/store"

	"github.com/micro/dev/model"

	"github.com/gosimple/slug"
	pb "github.com/micro/services/blog/posts/proto/posts"
	posts "github.com/micro/services/blog/posts/proto/posts"
	tags "github.com/micro/services/blog/tags/proto"
)

const (
	tagType = "post-tag"
)

type Post struct {
	ID      string   `json:"id"`
	Title   string   `json:"title"`
	Slug    string   `json:"slug"`
	Content string   `json:"content"`
	Created int64    `json:"created"`
	Updated int64    `json:"updated"`
	Tags    []string `json:"tags"`
}

type Posts struct {
	Tags tags.TagsService
	db   model.DB
}

func NewPosts(tagsService tags.TagsService) *Posts {
	createdIndex := model.ByEq("created")
	createdIndex.Reverse = true

	return &Posts{
		Tags: tagsService,
		db: model.NewDB(
			store.DefaultStore,
			"posts",
			model.Indexes(model.ByEq("id"), model.ByEq("slug"), createdIndex),
		),
	}
}

func (p *Posts) Save(ctx context.Context, req *posts.SaveRequest, rsp *posts.SaveResponse) error {
	if len(req.Id) == 0 {
		return errors.BadRequest("posts.save.input-check", "Id is missing")
	}

	// read by post
	posts := []Post{}
	err := p.db.List(model.Equals("id", req.Id), &posts)
	if err != nil {
		return errors.InternalServerError("posts.save.store-id-read", "Failed to read post by id: %v", err.Error())
	}
	postSlug := slug.Make(req.Title)
	// If no existing record is found, create a new one
	if len(posts) == 0 {
		post := &Post{
			ID:      req.Id,
			Title:   req.Title,
			Content: req.Content,
			Tags:    req.Tags,
			Slug:    postSlug,
			Created: time.Now().Unix(),
		}
		err := p.savePost(ctx, nil, post)
		if err != nil {
			return errors.InternalServerError("posts.save.post-save", "Failed to save new post: %v", err.Error())
		}
		return nil
	}
	oldPost := &posts[0]

	post := &Post{
		ID:      req.Id,
		Title:   oldPost.Title,
		Content: oldPost.Content,
		Slug:    oldPost.Slug,
		Tags:    oldPost.Tags,
		Created: oldPost.Created,
		Updated: time.Now().Unix(),
	}
	if len(req.Title) > 0 {
		post.Title = req.Title
		post.Slug = slug.Make(post.Title)
	}
	if len(req.Slug) > 0 {
		post.Slug = req.Slug
	}
	if len(req.Content) > 0 {
		post.Content = req.Content
	}
	if len(req.Tags) > 0 {
		// Handle the special case of deletion
		if len(req.Tags) == 0 && req.Tags[0] == "" {
			post.Tags = []string{}
		} else {
			post.Tags = req.Tags
		}
	}

	postsWithThisSlug := []Post{}
	err = p.db.List(model.Equals("slug", postSlug), &postsWithThisSlug)
	if err != nil {
		return errors.InternalServerError("posts.save.store-read", "Failed to read post by slug: %v", err.Error())
	}

	if len(postsWithThisSlug) > 0 {
		if oldPost.ID != postsWithThisSlug[0].ID {
			return errors.BadRequest("posts.save.slug-check", "An other post with this slug already exists")
		}
	}

	return p.savePost(ctx, oldPost, post)
}

func (p *Posts) savePost(ctx context.Context, oldPost, post *Post) error {
	err := p.db.Save(post)
	if err != nil {
		return err
	}
	if oldPost == nil {
		for _, tagName := range post.Tags {
			_, err := p.Tags.Add(ctx, &tags.AddRequest{
				ResourceID: post.ID,
				Type:       tagType,
				Title:      tagName,
			})
			if err != nil {
				return err
			}
		}
		return nil
	}
	return p.diffTags(ctx, post.ID, oldPost.Tags, post.Tags)
}

func (p *Posts) diffTags(ctx context.Context, parentID string, oldTagNames, newTagNames []string) error {
	oldTags := map[string]struct{}{}
	for _, v := range oldTagNames {
		oldTags[v] = struct{}{}
	}
	newTags := map[string]struct{}{}
	for _, v := range newTagNames {
		newTags[v] = struct{}{}
	}
	for i := range oldTags {
		_, stillThere := newTags[i]
		if !stillThere {
			_, err := p.Tags.Remove(ctx, &tags.RemoveRequest{
				ResourceID: parentID,
				Type:       tagType,
				Title:      i,
			})
			if err != nil {
				logger.Errorf("Error decreasing count for tag '%v' with type '%v' for parent '%v'", i, tagType, parentID)
			}
		}
	}
	for i := range newTags {
		_, newlyAdded := oldTags[i]
		if newlyAdded {
			_, err := p.Tags.Add(ctx, &tags.AddRequest{
				ResourceID: parentID,
				Type:       tagType,
				Title:      i,
			})
			if err != nil {
				logger.Errorf("Error increasing count for tag '%v' with type '%v' for parent '%v': %v", i, tagType, parentID, err)
			}
		}
	}
	return nil
}

func (p *Posts) Query(ctx context.Context, req *pb.QueryRequest, rsp *pb.QueryResponse) error {
	var q model.Query
	if len(req.Slug) > 0 {
		logger.Infof("Reading post by slug: %v", req.Slug)
		q = model.Equals("slug", req.Slug)
	} else if len(req.Id) > 0 {
		logger.Infof("Reading post by id: %v", req.Id)
		q = model.Equals("id", req.Id)
	} else {
		q := model.Equals("created", nil)
		var limit uint
		limit = 20
		if req.Limit > 0 {
			limit = uint(req.Limit)
		}
		q.Limit = limit
		q.Offset = req.Offset
		logger.Infof("Listing posts, offset: %v, limit: %v", req.Offset, limit)
	}

	posts := []Post{}
	err := p.db.List(q, &posts)
	if err != nil {
		return errors.BadRequest("posts.query.store-read", "Failed to read from store: %v", err.Error())
	}
	rsp.Posts = make([]*pb.Post, len(posts))
	for i, post := range posts {
		rsp.Posts[i] = &pb.Post{
			Id:      post.ID,
			Title:   post.Title,
			Slug:    post.Slug,
			Content: post.Content,
			Tags:    post.Tags,
		}
	}
	return nil
}

func (p *Posts) Delete(ctx context.Context, req *pb.DeleteRequest, rsp *pb.DeleteResponse) error {
	logger.Info("Received Post.Delete request")
	records, err := store.Read(fmt.Sprintf("%v:%v", idPrefix, req.Id))
	if err != nil && err != store.ErrNotFound {
		return err
	}
	if len(records) == 0 {
		return fmt.Errorf("Post with ID %v not found", req.Id)
	}
	post := &Post{}
	err = json.Unmarshal(records[0].Value, post)
	if err != nil {
		return err
	}

	// Delete by ID
	err = store.Delete(fmt.Sprintf("%v:%v", idPrefix, post.ID))
	if err != nil {
		return err
	}
	// Delete by slug
	err = store.Delete(fmt.Sprintf("%v:%v", slugPrefix, post.Slug))
	if err != nil {
		return err
	}
	// Delete by timeStamp
	return store.Delete(fmt.Sprintf("%v:%v", timeStampPrefix, post.CreateTimestamp))
}
