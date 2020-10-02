package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"time"

	"github.com/micro/go-micro/v3/errors"
	gostore "github.com/micro/go-micro/v3/store"
	"github.com/micro/micro/v3/service/logger"
	"github.com/micro/micro/v3/service/store"

	"github.com/gosimple/slug"
	pb "github.com/micro/services/blog/posts/proto/posts"
	posts "github.com/micro/services/blog/posts/proto/posts"
	tags "github.com/micro/services/blog/tags/proto"
)

const (
	tagType         = "post-tag"
	slugPrefix      = "slug"
	idPrefix        = "id"
	timeStampPrefix = "timestamp"
)

type Post struct {
	ID              string   `json:"id"`
	Title           string   `json:"title"`
	Slug            string   `json:"slug"`
	Content         string   `json:"content"`
	CreateTimestamp int64    `json:"create_timestamp"`
	UpdateTimestamp int64    `json:"update_timestamp"`
	TagNames        []string `json:"tagNames"`
}

type Posts struct {
	Tags tags.TagsService
}

func (p *Posts) Save(ctx context.Context, req *posts.SaveRequest, rsp *posts.SaveResponse) error {
	if len(req.Post.Id) == 0 || len(req.Post.Title) == 0 || len(req.Post.Content) == 0 {
		return errors.BadRequest("posts.save.input-check", "Id, title or content is missing")
	}

	// read by post
	records, err := store.Read(fmt.Sprintf("%v:%v", idPrefix, req.Post.Id))
	if err != nil && err != gostore.ErrNotFound {
		return errors.InternalServerError("posts.save.store-id-read", "Failed to read post by id: %v", err.Error())
	}
	postSlug := slug.Make(req.Post.Title)
	// If no existing record is found, create a new one
	if len(records) == 0 {
		post := &Post{
			ID:              req.Post.Id,
			Title:           req.Post.Title,
			Content:         req.Post.Content,
			TagNames:        req.Post.TagNames,
			Slug:            postSlug,
			CreateTimestamp: time.Now().Unix(),
		}
		err := p.savePost(ctx, nil, post)
		if err != nil {
			return errors.InternalServerError("posts.save.post-save", "Failed to save new post: %v", err.Error())
		}
		return nil
	}
	record := records[0]
	oldPost := &Post{}
	err = json.Unmarshal(record.Value, oldPost)
	if err != nil {
		return errors.InternalServerError("posts.save.unmarshal", "Failed to unmarshal old post: %v", err.Error())
	}
	post := &Post{
		ID:              req.Post.Id,
		Title:           req.Post.Title,
		Content:         req.Post.Content,
		Slug:            postSlug,
		TagNames:        req.Post.TagNames,
		CreateTimestamp: oldPost.CreateTimestamp,
		UpdateTimestamp: time.Now().Unix(),
	}

	// Check if slug exists
	recordsBySlug, err := store.Read(fmt.Sprintf("%v:%v", slugPrefix, postSlug))
	if err != nil && err != gostore.ErrNotFound {
		return errors.InternalServerError("posts.save.store-read", "Failed to read post by slug: %v", err.Error())
	}
	otherSlugPost := &Post{}
	err = json.Unmarshal(record.Value, otherSlugPost)
	if err != nil {
		return errors.InternalServerError("posts.save.slug-unmarshal", "Error unmarshaling other post with same slug: %v", err.Error())
	}
	if len(recordsBySlug) > 0 && oldPost.ID != otherSlugPost.ID {
		return errors.BadRequest("posts.save.slug-check", "An other post with this slug already exists")
	}

	return p.savePost(ctx, oldPost, post)
}

func (p *Posts) savePost(ctx context.Context, oldPost, post *Post) error {
	bytes, err := json.Marshal(post)
	if err != nil {
		return err
	}

	err = store.Write(&gostore.Record{
		Key:   fmt.Sprintf("%v:%v", idPrefix, post.ID),
		Value: bytes,
	})
	if err != nil {
		return err
	}
	// Delete old slug index if the slug has changed
	if oldPost != nil && oldPost.Slug != post.Slug {
		err = store.Delete(fmt.Sprintf("%v:%v", slugPrefix, post.Slug))
		if err != nil {
			return err
		}
	}
	err = store.Write(&gostore.Record{
		Key:   fmt.Sprintf("%v:%v", slugPrefix, post.Slug),
		Value: bytes,
	})
	if err != nil {
		return err
	}
	err = store.Write(&gostore.Record{
		Key:   fmt.Sprintf("%v:%v", timeStampPrefix, math.MaxInt64-post.CreateTimestamp),
		Value: bytes,
	})
	if err != nil {
		return err
	}
	if oldPost == nil {
		for _, tagName := range post.TagNames {
			_, err := p.Tags.IncreaseCount(ctx, &tags.IncreaseCountRequest{
				ParentID: post.ID,
				Type:     tagType,
				Title:    tagName,
			})
			if err != nil {
				return err
			}
		}
		return nil
	}
	return p.diffTags(ctx, post.ID, oldPost.TagNames, post.TagNames)
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
			_, err := p.Tags.DecreaseCount(ctx, &tags.DecreaseCountRequest{
				ParentID: parentID,
				Type:     tagType,
				Title:    i,
			})
			if err != nil {
				logger.Errorf("Error decreasing count for tag '%v' with type '%v' for parent '%v'", i, tagType, parentID)
			}
		}
	}
	for i := range newTags {
		_, newlyAdded := oldTags[i]
		if newlyAdded {
			_, err := p.Tags.IncreaseCount(ctx, &tags.IncreaseCountRequest{
				ParentID: parentID,
				Type:     tagType,
				Title:    i,
			})
			if err != nil {
				logger.Errorf("Error increasing count for tag '%v' with type '%v' for parent '%v'", i, tagType, parentID)
			}
		}
	}
	return nil
}

func (p *Posts) Query(ctx context.Context, req *pb.QueryRequest, rsp *pb.QueryResponse) error {
	var records []*gostore.Record
	var err error
	if len(req.Slug) > 0 {
		key := fmt.Sprintf("%v:%v", slugPrefix, req.Slug)
		logger.Infof("Reading post by slug: %v", req.Slug)
		records, err = store.Read("", store.Prefix(key))
	} else {
		key := fmt.Sprintf("%v:", timeStampPrefix)
		var limit uint
		limit = 20
		if req.Limit > 0 {
			limit = uint(req.Limit)
		}
		logger.Infof("Listing posts, offset: %v, limit: %v", req.Offset, limit)
		records, err = store.Read("", store.Prefix(key),
			store.Offset(uint(req.Offset)),
			store.Limit(limit))
	}

	if err != nil {
		return errors.BadRequest("posts.query.store-read", "Failed to read from store: %v", err.Error())
	}
	rsp.Posts = make([]*pb.Post, len(records))
	for i, record := range records {
		postRecord := &Post{}
		err := json.Unmarshal(record.Value, postRecord)
		if err != nil {
			return errors.InternalServerError("posts.save.unmarshal", "Failed to unmarshal old post: %v", err.Error())
		}
		rsp.Posts[i] = &pb.Post{
			Id:       postRecord.ID,
			Title:    postRecord.Title,
			Slug:     postRecord.Slug,
			Content:  postRecord.Content,
			TagNames: postRecord.TagNames,
		}
	}
	return nil
}

func (p *Posts) Delete(ctx context.Context, req *pb.DeleteRequest, rsp *pb.DeleteResponse) error {
	logger.Info("Received Post.Delete request")
	records, err := store.Read(fmt.Sprintf("%v:%v", idPrefix, req.Id))
	if err != nil && err != gostore.ErrNotFound {
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
