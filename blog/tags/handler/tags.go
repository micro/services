package handler

import (
	"context"
	"encoding/json"

	"fmt"

	"github.com/gosimple/slug"
	"github.com/micro/dev/model"
	"github.com/micro/micro/v3/service/errors"
	"github.com/micro/micro/v3/service/logger"
	"github.com/micro/micro/v3/service/store"
	proto "github.com/micro/services/blog/tags/proto"
)

const (
	resourcePrefix = "byResource"
	tagCountPrefix = "tagCount"
	childrenByTag  = "childrenByTag"
)

type Tags struct {
	db model.Model
}

func NewTags() *Tags {
	slugIndex := model.ByEquality("slug")
	slugIndex.Order.Type = model.OrderTypeUnordered
	return &Tags{
		db: model.New(
			store.DefaultStore,
			"tags",
			model.Indexes(model.ByEquality("type")),
			&model.ModelOptions{
				IdIndex: slugIndex,
				Debug:   false,
			},
		),
	}
}

func (t *Tags) Add(ctx context.Context, req *proto.AddRequest, rsp *proto.AddResponse) error {
	if len(req.ResourceID) == 0 || len(req.Type) == 0 {
		return errors.BadRequest("tags.increasecount.input-check", "resource id and type is required")
	}

	tags := []*proto.Tag{}
	tagSlug := slug.Make(req.GetTitle())
	q := model.Equals("slug", tagSlug)
	q.Order.Type = model.OrderTypeUnordered
	err := t.db.List(q, &tags)
	if err != nil {
		return err
	}

	var tag *proto.Tag
	// If no existing record is found, create a new one
	if len(tags) == 0 {
		tag = &proto.Tag{
			Title: req.GetTitle(),
			Type:  req.Type,
			Slug:  tagSlug,
		}
	} else {
		tag = tags[0]
	}

	// increase tag count
	err = store.Write(&store.Record{
		Key:   fmt.Sprintf("%v:%v:%v", tagCountPrefix, tag.Slug, req.GetResourceID()),
		Value: nil,
	})
	if err != nil {
		return err
	}

	oldTagCount := tag.Count
	// get tag count
	recs, err := store.List(store.Prefix(fmt.Sprintf("%v:%v", tagCountPrefix, tag.Slug)), store.Limit(1000))
	if err != nil {
		return err
	}

	tag.Count = int64(len(recs))
	if tag.Count == oldTagCount {
		return fmt.Errorf("Tag count for tag %v is unchanged, was: %v, now: %v", tagSlug, oldTagCount, tag.Count)
	}
	tagJSON, err := json.Marshal(tag)
	if err != nil {
		return err
	}
	err = store.Write(&store.Record{
		Key:   fmt.Sprintf("%v:%v:%v", resourcePrefix, req.GetResourceID(), tag.Slug),
		Value: tagJSON,
	})
	if err != nil {
		return err
	}
	return t.saveTag(tag)
}

func (t *Tags) saveTag(tag *proto.Tag) error {
	tag.Slug = slug.Make(tag.Title)
	return t.db.Save(tag)
}

func (t *Tags) Remove(ctx context.Context, req *proto.RemoveRequest, rsp *proto.RemoveResponse) error {
	if len(req.ResourceID) == 0 || len(req.Type) == 0 {
		return errors.BadRequest("tags.decreaseecount.input-check", "resource id and type is required")
	}

	tagSlug := slug.Make(req.GetTitle())
	resourceKey := fmt.Sprintf("%v:%v:%v", resourcePrefix, req.GetResourceID(), tagSlug)

	// read by resource ID + slug, the record is identical in boths places anyway
	records, err := store.Read(resourceKey)
	if err != nil && err != store.ErrNotFound {
		return err
	}

	// If no existing record is found, there is nothing to decrease
	if len(records) == 0 {
		// return error?
		return nil
	}
	record := records[0]
	tag := &proto.Tag{}
	err = json.Unmarshal(record.Value, tag)
	if err != nil {
		return err
	}

	// decrease tag count
	err = store.Delete(fmt.Sprintf("%v:%v:%v", tagCountPrefix, tag.Slug, req.GetResourceID()))
	if err != nil {
		return err
	}

	// get tag count
	recs, err := store.List(store.Prefix(fmt.Sprintf("%v:%v", tagCountPrefix, tag.Slug)), store.Limit(1000))
	if err != nil {
		return err
	}
	tag.Count = int64(len(recs))
	return t.saveTag(tag)
}

func (t *Tags) List(ctx context.Context, req *proto.ListRequest, rsp *proto.ListResponse) error {
	logger.Info("Received Tags.List request")

	// unfortunately there is a mixing of manual indexes
	// and model here because model does not yet support
	// many to many relations
	key := ""
	var q model.Query
	if len(req.ResourceID) > 0 {
		key = fmt.Sprintf("%v:%v", resourcePrefix, req.ResourceID)
	} else if len(req.Type) > 0 {
		q = model.Equals("type", req.Type)
	} else {
		return errors.BadRequest("tags.list.input-check", "resource id or type is required")
	}

	if q.Type != "" {
		tags := []proto.Tag{}
		err := t.db.List(q, &tags)
		if err != nil {
			return err
		}
		rsp.Tags = make([]*proto.Tag, len(tags))
		for i, tag := range tags {
			rsp.Tags[i] = &proto.Tag{
				Title: tag.Title,
				Type:  tag.Type,
				Slug:  tag.Slug,
				Count: tag.Count,
			}
		}
		return nil
	}
	records, err := store.Read("", store.Prefix(key))
	if err != nil {
		return err
	}

	rsp.Tags = make([]*proto.Tag, len(records))
	for i, record := range records {
		tagRecord := &proto.Tag{}
		err := json.Unmarshal(record.Value, tagRecord)
		if err != nil {
			return err
		}
		rsp.Tags[i] = &proto.Tag{
			Title: tagRecord.Title,
			Type:  tagRecord.Type,
			Slug:  tagRecord.Slug,
			Count: tagRecord.Count,
		}
	}

	return nil
}

func (t *Tags) Update(ctx context.Context, req *proto.UpdateRequest, rsp *proto.UpdateResponse) error {
	if len(req.Title) == 0 || len(req.Type) == 0 {
		return errors.BadRequest("tags.update.input-check", "title and type is required")
	}

	tagSlug := slug.Make(req.GetTitle())
	tags := []proto.Tag{}
	q := model.Equals("slug", tagSlug)
	q.Order.Type = model.OrderTypeUnordered
	err := t.db.List(q, &tags)
	if err != nil {
		return err
	}
	if len(tags) == 0 {
		return errors.BadRequest("tags.update.input-check", "Tag not found")
	}
	tags[0].Title = req.Title
	return t.saveTag(&tags[0])
}
