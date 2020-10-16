package handler

import (
	"context"
	"encoding/json"

	"fmt"

	"github.com/gosimple/slug"
	"github.com/micro/go-micro/v3/errors"
	"github.com/micro/micro/v3/service/logger"
	"github.com/micro/micro/v3/service/store"
	pb "github.com/micro/services/blog/tags/proto"
)

const (
	slugPrefix     = "bySlug"
	resourcePrefix = "byResource"
	typePrefix     = "byType"
	tagCountPrefix = "tagCount"
	childrenByTag  = "childrenByTag"
)

type Tag struct {
	Title string `json:"title"`
	Slug  string `json:"slug"`
	Type  string `json:"type"`
	Count int64  `json:"count"`
}

type Tags struct{}

func (t *Tags) Add(ctx context.Context, req *pb.AddRequest, rsp *pb.AddResponse) error {
	if len(req.ResourceID) == 0 || len(req.Type) == 0 {
		return errors.BadRequest("tags.increasecount.input-check", "resource id and type is required")
	}

	tagSlug := slug.Make(req.GetTitle())
	key := fmt.Sprintf("%v:%v", slugPrefix, tagSlug)

	// read by resource ID + slug, the record is identical in boths places anyway
	records, err := store.Read(key)
	if err != nil && err != store.ErrNotFound {
		return err
	}

	var tag *Tag
	// If no existing record is found, create a new one
	if len(records) == 0 {
		tag = &Tag{
			Title: req.GetTitle(),
			Type:  req.Type,
			Slug:  tagSlug,
		}
	} else {
		record := records[0]
		tag = &Tag{}
		err = json.Unmarshal(record.Value, tag)
		if err != nil {
			return err
		}
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

func (t *Tags) saveTag(tag *Tag) error {
	tagSlug := slug.Make(tag.Title)

	key := fmt.Sprintf("%v:%v", slugPrefix, tagSlug)
	typeKey := fmt.Sprintf("%v:%v:%v", typePrefix, tag.Type, tagSlug)

	bytes, err := json.Marshal(tag)
	if err != nil {
		return err
	}

	// write resourceId:slug to enable prefix listing based on type
	err = store.Write(&store.Record{
		Key:   key,
		Value: bytes,
	})
	if err != nil {
		return err
	}
	return store.Write(&store.Record{
		Key:   typeKey,
		Value: bytes,
	})
}

func (t *Tags) Remove(ctx context.Context, req *pb.RemoveRequest, rsp *pb.RemoveResponse) error {
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
	tag := &Tag{}
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

func (t *Tags) List(ctx context.Context, req *pb.ListRequest, rsp *pb.ListResponse) error {
	logger.Info("Received Tags.List request")
	key := ""
	if len(req.ResourceID) > 0 {
		key = fmt.Sprintf("%v:%v", resourcePrefix, req.ResourceID)
	} else if len(req.Type) > 0 {
		key = fmt.Sprintf("%v:%v", typePrefix, req.Type)
	} else {
		return errors.BadRequest("tags.list.input-check", "resource id or type is required")
	}

	records, err := store.Read("", store.Prefix(key))
	if err != nil {
		return err
	}

	rsp.Tags = make([]*pb.Tag, len(records))
	for i, record := range records {
		tagRecord := &Tag{}
		err := json.Unmarshal(record.Value, tagRecord)
		if err != nil {
			return err
		}
		rsp.Tags[i] = &pb.Tag{
			Title: tagRecord.Title,
			Type:  tagRecord.Type,
			Slug:  tagRecord.Slug,
			Count: tagRecord.Count,
		}
	}

	return nil
}

func (t *Tags) Update(ctx context.Context, req *pb.UpdateRequest, rsp *pb.UpdateResponse) error {
	if len(req.Title) == 0 || len(req.Type) == 0 {
		return errors.BadRequest("tags.update.input-check", "title and type is required")
	}

	tagSlug := slug.Make(req.GetTitle())
	resourceID := fmt.Sprintf("%v:%v", slugPrefix, tagSlug)

	// read by resource ID + slug, the record is identical in boths places anyway
	records, err := store.Read(resourceID)
	if err != nil {
		return err
	}

	if len(records) == 0 {
		return fmt.Errorf("Tag with slug '%v' not found, nothing to update", tagSlug)
	}
	record := records[0]
	tag := &Tag{}
	err = json.Unmarshal(record.Value, tag)
	if err != nil {
		return err
	}
	tag.Title = req.Title
	return t.saveTag(tag)
}
