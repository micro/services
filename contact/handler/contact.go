package handler

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"github.com/micro/micro/v3/service/errors"
	"github.com/micro/micro/v3/service/logger"
	pauth "github.com/micro/services/pkg/auth"
	adminpb "github.com/micro/services/pkg/service/proto"
	"github.com/micro/services/pkg/tenant"

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
		Addresses:    req.Addresses,
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
		Addresses:    req.Addresses,
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

func (c *contact) DeleteData(ctx context.Context, request *adminpb.DeleteDataRequest, response *adminpb.DeleteDataResponse) error {
	method := "admin.DeleteData"
	_, err := pauth.VerifyMicroAdmin(ctx, method)
	if err != nil {
		return err
	}

	if len(request.TenantId) < 10 { // deliberate length check so we don't delete all the things
		return errors.BadRequest(method, "Missing tenant ID")
	}

	split := strings.Split(request.TenantId, "/")
	tctx := tenant.NewContext(split[1], split[0], split[1])
	// load all keys
	keys := []string{}
	offset := uint(0)
	for {
		res, err := c.contact.List(tctx, offset, 100)
		if err != nil && !strings.Contains(err.Error(), "not found") {
			return err
		}
		for _, r := range res {
			keys = append(keys, r.Id)
		}
		if len(res) < 100 {
			break
		}
		offset += 100
	}
	for _, k := range keys {
		if err := c.contact.Delete(tctx, k); err != nil {
			return err
		}
	}

	logger.Infof("Deleted %d keys for %s", len(keys), request.TenantId)
	return nil
}
