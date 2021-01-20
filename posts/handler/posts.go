package handler

import (
	"context"
	"time"

	"github.com/micro/micro/v3/service/errors"
	"github.com/micro/micro/v3/service/logger"
	"github.com/micro/micro/v3/service/store"

	"github.com/micro/dev/model"

	"github.com/gosimple/slug"
	proto "github.com/micro/services/posts/proto"
	tags "github.com/micro/services/tags/proto"
)

const (
	tagType = "post-tag"
)

type Posts struct {
	Tags tags.TagsService
	db   model.Model
}

func NewPosts(tagsService tags.TagsService) *Posts {
	createdIndex := model.ByEquality("created")
	createdIndex.Order.Type = model.OrderTypeDesc

	return &Posts{
		Tags: tagsService,
		db: model.New(
			store.DefaultStore,
			"posts",
			model.Indexes(model.ByEquality("slug"), createdIndex),
			&model.ModelOptions{
				Debug: false,
			},
		),
	}
}

func (p *Posts) Save(ctx context.Context, req *proto.SaveRequest, rsp *proto.SaveResponse) error {
	if len(req.Id) == 0 {
		return errors.BadRequest("proto.save.input-check", "Id is missing")
	}

	// read by post
	posts := []*proto.Post{}
	q := model.Equals("id", req.Id)
	q.Order.Type = model.OrderTypeUnordered
	err := p.db.List(q, &posts)
	if err != nil {
		return errors.InternalServerError("proto.save.store-id-read", "Failed to read post by id: %v", err.Error())
	}
	postSlug := slug.Make(req.Title)
	// If no existing record is found, create a new one
	if len(posts) == 0 {
		post := &proto.Post{
			Id:       req.Id,
			Title:    req.Title,
			Content:  req.Content,
			Tags:     req.Tags,
			Slug:     postSlug,
			Created:  time.Now().Unix(),
			Metadata: req.Metadata,
			Image:    req.Image,
		}
		err := p.savePost(ctx, nil, post)
		if err != nil {
			return errors.InternalServerError("proto.save.post-save", "Failed to save new post: %v", err.Error())
		}
		return nil
	}
	oldPost := posts[0]

	post := &proto.Post{
		Id:       req.Id,
		Title:    oldPost.Title,
		Content:  oldPost.Content,
		Slug:     oldPost.Slug,
		Tags:     oldPost.Tags,
		Created:  oldPost.Created,
		Updated:  time.Now().Unix(),
		Metadata: req.Metadata,
		Image:    req.Image,
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

	postsWithThisSlug := []*proto.Post{}
	err = p.db.List(model.Equals("slug", postSlug), &postsWithThisSlug)
	if err != nil {
		return errors.InternalServerError("proto.save.store-read", "Failed to read post by slug: %v", err.Error())
	}

	if len(postsWithThisSlug) > 0 {
		if oldPost.Id != postsWithThisSlug[0].Id {
			return errors.BadRequest("proto.save.slug-check", "An other post with this slug already exists")
		}
	}

	return p.savePost(ctx, oldPost, post)
}

func (p *Posts) savePost(ctx context.Context, oldPost, post *proto.Post) error {
	err := p.db.Save(post)
	if err != nil {
		return err
	}
	if oldPost == nil {
		for _, tagName := range post.Tags {
			_, err := p.Tags.Add(ctx, &tags.AddRequest{
				ResourceID: post.Id,
				Type:       tagType,
				Title:      tagName,
			})
			if err != nil {
				return err
			}
		}
		return nil
	}
	return p.diffTags(ctx, post.Id, oldPost.Tags, post.Tags)
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

func (p *Posts) Query(ctx context.Context, req *proto.QueryRequest, rsp *proto.QueryResponse) error {
	var q model.Query
	if len(req.Slug) > 0 {
		logger.Infof("Reading post by slug: %v", req.Slug)
		q = model.Equals("slug", req.Slug)
	} else if len(req.Id) > 0 {
		logger.Infof("Reading post by id: %v", req.Id)
		q = model.Equals("id", req.Id)
		q.Order.Type = model.OrderTypeUnordered
	} else {
		q = model.Equals("created", nil)
		q.Order.Type = model.OrderTypeDesc
		var limit uint
		limit = 20
		if req.Limit > 0 {
			limit = uint(req.Limit)
		}
		q.Limit = int64(limit)
		q.Offset = req.Offset
		logger.Infof("Listing posts, offset: %v, limit: %v", req.Offset, limit)
	}

	return p.db.List(q, &rsp.Posts)
}

func (p *Posts) Delete(ctx context.Context, req *proto.DeleteRequest, rsp *proto.DeleteResponse) error {
	logger.Info("Received Post.Delete request")
	q := model.Equals("id", req.Id)
	q.Order.Type = model.OrderTypeUnordered
	return p.db.Delete(q)
}
