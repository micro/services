package domain

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/micro/micro/v3/service/store"
	"github.com/pkg/errors"

	pb "github.com/micro/services/contact/proto"
	"github.com/micro/services/pkg/tenant"
)

type Contact interface {
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

// contactIdPrefix return the contact prefix of the store key
func contactIdPrefix(ctx context.Context) string {
	tenantId, ok := tenant.FromContext(ctx)
	if !ok {
		tenantId = "micro"
	}

	return fmt.Sprintf("contact/%s/id/", tenantId)
}

// contactIdPrefix return the store key of one contact
func contactIdKey(ctx context.Context, id string) string {
	return fmt.Sprintf("%s%s", contactIdPrefix(ctx), id)
}

// Create a contact
func (c *contact) Create(ctx context.Context, info *pb.ContactInfo) error {
	info.CreatedAt = time.Now().Format(time.RFC3339)
	info.UpdatedAt = time.Now().Format(time.RFC3339)

	val, err := json.Marshal(info)
	if err != nil {
		return err
	}

	return store.Write(&store.Record{
		Key:   contactIdKey(ctx, info.Id),
		Value: val,
	})
}

// Update contact information by id
func (c *contact) Update(ctx context.Context, id string, info *pb.ContactInfo) error {
	info.UpdatedAt = time.Now().Format(time.RFC3339)

	val, err := json.Marshal(info)
	if err != nil {
		return err
	}

	return store.Write(&store.Record{
		Key:   contactIdKey(ctx, id),
		Value: val,
	})
}

// Read one contact by id
func (c *contact) Read(ctx context.Context, id string) (*pb.ContactInfo, error) {
	records, err := c.store.Read(contactIdKey(ctx, id))
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

// Delete one contact by id
func (c *contact) Delete(ctx context.Context, id string) error {
	return c.store.Delete(contactIdKey(ctx, id))
}

// List contacts by offset and limit
func (c *contact) List(ctx context.Context, offset, limit uint) (result []*pb.ContactInfo, err error) {
	records, err := c.store.Read(contactIdPrefix(ctx),
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
