package handler

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"github.com/micro/micro/v3/service/errors"

	"github.com/micro/services/contact/domain"
	pb "github.com/micro/services/contact/proto"
)

type contact struct {
	contact domain.Contact
}

func NewContact(c domain.Contact) *contact {
	return &contact{
		contact: c,
	}
}

// Create a contact
func (c *contact) Create(ctx context.Context, req *pb.CreateRequest, rsp *pb.CreateResponse) error {
	req.Name = strings.TrimSpace(req.Name)
	if len(req.Name) == 0 {
		return errors.BadRequest("contact.create", "contact name is required")
	}

	uid, err := uuid.NewUUID()
	if err != nil {
		return errors.InternalServerError("contact.create", "generate contact id error: %v", err)
	}

	info := &pb.ContactInfo{
		Id:           uid.String(),
		Name:         req.Name,
		Phones:       req.Phones,
		Emails:       req.Emails,
		Links:        req.Links,
		Birthday:     req.Birthday,
		Locations:    req.Locations,
		SocialMedias: req.SocialMedias,
		Note:         req.Note,
	}

	if err := c.contact.Create(ctx, info); err != nil {
		return errors.InternalServerError("contact.create", "create contact error: %v", err)
	}

	rsp.Contact = info

	return nil
}

// Update information of the contact submitted
func (c *contact) Update(ctx context.Context, req *pb.UpdateRequest, rsp *pb.UpdateResponse) error {
	req.Name = strings.TrimSpace(req.Name)
	if len(req.Name) == 0 {
		return errors.BadRequest("contact.create", "contact name is required")
	}

	old, err := c.contact.Read(ctx, req.Id)
	if err != nil {
		return errors.InternalServerError("contact.update", "get contact info error: %v", err)
	}

	info := &pb.ContactInfo{
		Id:           req.Id,
		Name:         req.Name,
		Phones:       req.Phones,
		Emails:       req.Emails,
		Links:        req.Links,
		Birthday:     req.Birthday,
		Locations:    req.Locations,
		SocialMedias: req.SocialMedias,
		Note:         req.Note,
		CreatedAt:    old.CreatedAt,
	}

	if err := c.contact.Update(ctx, req.Id, info); err != nil {
		return errors.InternalServerError("contact.update", "update contact error: %v", err)
	}

	rsp.Contact = info

	return nil
}

// Read a contact by id
func (c *contact) Read(ctx context.Context, req *pb.ReadRequest, rsp *pb.ReadResponse) error {
	info, err := c.contact.Read(ctx, req.Id)
	if err != nil {
		return errors.InternalServerError("contact.read", "get contact info error: %v", err)
	}

	rsp.Contact = info

	return nil
}

// Delete contact by id
func (c *contact) Delete(ctx context.Context, req *pb.DeleteRequest, _ *pb.DeleteResponse) error {
	err := c.contact.Delete(ctx, req.Id)
	if err != nil {
		return errors.InternalServerError("contact.delete", "delete contact error: %v", err)
	}

	return nil
}

// List contacts with offset and limit
func (c *contact) List(ctx context.Context, req *pb.ListRequest, rsp *pb.ListResponse) error {
	if req.Limit == 0 {
		req.Limit = 30
	}

	list, err := c.contact.List(ctx, uint(req.Offset), uint(req.Limit))
	if err != nil {
		return errors.InternalServerError("contact.list", "get contact info error: %v", err)
	}

	rsp.Contacts = list

	return nil
}
