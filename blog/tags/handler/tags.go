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
	slugPrefix   = "bySlug"
	parentPrefix = "byParent"
	typePrefix   = "byType"
)

type Tag struct {
	Title string `json:"title"`
	Slug  string `json:"slug"`
	Type  string `json:"type"`
	Count int64  `json:"count"`
}

type Tags struct{}

func (t *Tags) Add(ctx context.Context, req *pb.AddRequest, rsp *pb.AddResponse) error {
	if len(req.ParentID) == 0 || len(req.Type) == 0 {
		return errors.BadRequest("tags.increasecount.input-check", "parent id and type is required")
	}

	tagSlug := slug.Make(req.GetTitle())
	parentID := fmt.Sprintf("%v:%v:%v", parentPrefix, req.GetParentID(), tagSlug)

	// read by parent ID + slug, the record is identical in boths places anyway
	records, err := store.Read(parentID)
	if err != nil && err != store.ErrNotFound {
		return err
	}

	// If no existing record is found, create a new one
	if len(records) == 0 {
		tag := &Tag{
			Title: req.GetTitle(),
			Type:  req.Type,
			Slug:  tagSlug,
			Count: 1,
		}
		return t.saveTag(tag)
	}
	record := records[0]
	tag := &Tag{}
	err = json.Unmarshal(record.Value, tag)
	if err != nil {
		return err
	}
	recs, err := store.List(store.Prefix(fmt.Sprintf("%v:%v", parentPrefix, req.ParentID)), store.Limit(1000))
	if err != nil {
		return err
	}
	tag.Count = int64(len(recs))
	tagJSON, err := json.Marshal(tag)
	if err != nil {
		return err
	}
	err = store.Write(&store.Record{
		Key:   fmt.Sprintf("%v:%v:%v", parentPrefix, parentID, tag.Slug),
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

	// write parentId:slug to enable prefix listing based on type
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
	if len(req.ParentID) == 0 || len(req.Type) == 0 {
		return errors.BadRequest("tags.decreaseecount.input-check", "parent id and type is required")
	}

	tagSlug := slug.Make(req.GetTitle())
	parentKey := fmt.Sprintf("%v:%v:%v", parentPrefix, req.GetParentID(), tagSlug)

	// read by parent ID + slug, the record is identical in boths places anyway
	records, err := store.Read(parentKey)
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
	err = store.Delete(fmt.Sprintf("%v:%v:%v", parentPrefix, req.GetParentID(), tag.Slug))
	if err != nil {
		return err
	}
	recs, err := store.List(store.Prefix(fmt.Sprintf("%v:%v", parentPrefix, req.ParentID)), store.Limit(1000))
	if err != nil {
		return err
	}
	tag.Count = int64(len(recs))
	return t.saveTag(tag)
}

func (t *Tags) List(ctx context.Context, req *pb.ListRequest, rsp *pb.ListResponse) error {
	logger.Info("Received Tags.List request")
	key := ""
	if len(req.ParentID) > 0 {
		key = fmt.Sprintf("%v:%v", parentPrefix, req.ParentID)
	} else if len(req.Type) > 0 {
		key = fmt.Sprintf("%v:%v", typePrefix, req.Type)
	} else {
		return errors.BadRequest("tags.list.input-check", "parent id or type is required")
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
	parentID := fmt.Sprintf("%v:%v", slugPrefix, tagSlug)

	// read by parent ID + slug, the record is identical in boths places anyway
	records, err := store.Read(parentID)
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
