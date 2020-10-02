package handler

import (
	"context"
	"encoding/json"

	"fmt"

	"github.com/gosimple/slug"
	"github.com/micro/go-micro/v3/errors"
	gostore "github.com/micro/go-micro/v3/store"
	"github.com/micro/micro/v3/service/logger"
	"github.com/micro/micro/v3/service/store"
	pb "github.com/micro/services/blog/tags/proto"
)

const (
	parentPrefix = "parent"
	typePrefix   = "type"
)

type Tag struct {
	ParentID string `json:"parentID"`
	Title    string `json:"title"`
	Slug     string `json:"slug"`
	Type     string `json:"type"`
	Count    int64  `json:"count"`
}

type Tags struct{}

func (t *Tags) IncreaseCount(ctx context.Context, req *pb.IncreaseCountRequest, rsp *pb.IncreaseCountResponse) error {
	if len(req.ParentID) == 0 || len(req.Type) == 0 {
		return errors.BadRequest("tags.increasecount.input-check", "parent id and type is required")
	}

	tagSlug := slug.Make(req.GetTitle())
	parentID := fmt.Sprintf("%v:%v:%v", parentPrefix, req.GetParentID(), tagSlug)

	// read by parent ID + slug, the record is identical in boths places anyway
	records, err := store.Read(parentID)
	if err != nil && err != gostore.ErrNotFound {
		return err
	}

	// If no existing record is found, create a new one
	if len(records) == 0 {
		tag := &Tag{
			ParentID: req.GetParentID(),
			Title:    req.GetTitle(),
			Type:     req.Type,
			Slug:     tagSlug,
			Count:    1,
		}
		return t.saveTag(tag)
	}
	record := records[0]
	tag := &Tag{}
	err = json.Unmarshal(record.Value, tag)
	if err != nil {
		return err
	}
	tag.Count++
	return t.saveTag(tag)
}

func (t *Tags) saveTag(tag *Tag) error {
	tagSlug := slug.Make(tag.Title)

	parentID := fmt.Sprintf("%v:%v:%v", parentPrefix, tag.ParentID, tagSlug)
	typeID := fmt.Sprintf("%v:%v:%v", typePrefix, tag.Type, tagSlug)

	bytes, err := json.Marshal(tag)
	if err != nil {
		return err
	}

	// write parentId:slug to enable prefix listing based on parent
	err = store.Write(&gostore.Record{
		Key:   parentID,
		Value: bytes,
	})
	if err != nil {
		return err
	}

	// write type:slug to enable prefix listing based on parent
	return store.Write(&gostore.Record{
		Key:   typeID,
		Value: bytes,
	})
}

func (t *Tags) DecreaseCount(ctx context.Context, req *pb.DecreaseCountRequest, rsp *pb.DecreaseCountResponse) error {
	if len(req.ParentID) == 0 || len(req.Type) == 0 {
		return errors.BadRequest("tags.decreaseecount.input-check", "parent id and type is required")
	}

	tagSlug := slug.Make(req.GetTitle())
	parentID := fmt.Sprintf("%v:%v:%v", parentPrefix, req.GetParentID(), tagSlug)

	// read by parent ID + slug, the record is identical in boths places anyway
	records, err := store.Read(parentID)
	if err != nil && err != gostore.ErrNotFound {
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
	if tag.Count == 0 {
		// return error?
		return nil
	}
	tag.Count--
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
			ParentID: tagRecord.ParentID,
			Title:    tagRecord.Title,
			Type:     tagRecord.Type,
			Slug:     tagRecord.Slug,
			Count:    tagRecord.Count,
		}
	}
	return nil
}

func (t *Tags) Update(ctx context.Context, req *pb.UpdateRequest, rsp *pb.UpdateResponse) error {
	if len(req.ParentID) == 0 || len(req.Type) == 0 {
		return errors.BadRequest("tags.update.input-check", "parent id and type is required")
	}

	tagSlug := slug.Make(req.GetTitle())
	parentID := fmt.Sprintf("%v:%v:%v", parentPrefix, req.GetParentID(), tagSlug)

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
