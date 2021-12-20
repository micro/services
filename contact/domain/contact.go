package domain

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/micro/micro/v3/service/store"
	"github.com/pkg/errors"

	pb "github.com/micro/services/contact/proto"
)

type ContactIface interface {
	Create(ctx context.Context, info *pb.ContactInfo) error
	Update(ctx context.Context, id string, info *pb.ContactInfo) error
	Read(ctx context.Context, id string) (result *pb.ContactInfo, err error)
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, offset, limit uint) (result []*pb.ContactInfo, err error)
}

type contact struct {
	store store.Store
}

func NewContactDomain(s store.Store) *contact {
	return &contact{
		store: s,
	}
}

func contactIdPrefix() string {
	return "contact/id/"
}

func contactIdKey(id string) string {
	return fmt.Sprintf("%s%s", contactIdPrefix(), id)
}

func (c *contact) Create(ctx context.Context, info *pb.ContactInfo) error {
	info.CreatedAt = time.Now().Unix()
	info.UpdatedAt = time.Now().Unix()

	val, err := json.Marshal(info)
	if err != nil {
		return err
	}

	return store.Write(&store.Record{
		Key:   contactIdKey(info.Id),
		Value: val,
	})
}

func (c *contact) Update(ctx context.Context, id string, info *pb.ContactInfo) error {
	info.UpdatedAt = time.Now().Unix()

	val, err := json.Marshal(info)
	if err != nil {
		return err
	}

	return store.Write(&store.Record{
		Key:   contactIdKey(id),
		Value: val,
	})
}

func (c *contact) Read(ctx context.Context, id string) (*pb.ContactInfo, error) {
	records, err := c.store.Read(contactIdKey(id))
	if err != nil {
		return nil, err
	}

	if len(records) == 0 {
		return nil, errors.New("not found")
	}

	info := &pb.ContactInfo{}
	if err := json.Unmarshal(records[0].Value, info); err != nil {
		return nil, err
	}

	return info, err
}

func (c *contact) Delete(ctx context.Context, id string) error {
	return c.store.Delete(contactIdKey(id))
}

func (c *contact) List(ctx context.Context, offset, limit uint) (result []*pb.ContactInfo, err error) {
	records, err := c.store.Read(contactIdPrefix(),
		store.ReadPrefix(),
		store.ReadOffset(offset),
		store.ReadLimit(limit))
	if err != nil {
		return nil, err
	}

	if len(records) == 0 {
		return nil, errors.New("not found")
	}

	for _, rec := range records {
		cinfo := &pb.ContactInfo{}
		json.Unmarshal(rec.Value, cinfo)
		result = append(result, cinfo)
	}

	return result, err
}
