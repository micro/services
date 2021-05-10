package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"regexp"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/micro/micro/v3/service/auth"
	"github.com/micro/micro/v3/service/errors"
	"github.com/micro/micro/v3/service/logger"
	"github.com/micro/micro/v3/service/store"
	pb "github.com/micro/services/invites/proto"
	"github.com/micro/services/pkg/tenant"
)

var (
	ErrMissingID              = errors.BadRequest("MISSING_ID", "Missing ID")
	ErrMissingGroupID         = errors.BadRequest("MISSING_GROUP_ID", "Missing GroupID")
	ErrInvalidEmail           = errors.BadRequest("INVALID_EMAIL", "The email provided was invalid")
	ErrMissingEmail           = errors.BadRequest("MISSING_EMAIL", "Missing Email")
	ErrMissingIDAndCode       = errors.BadRequest("ID_OR_CODE_REQUIRED", "An email address code is required to read an invite")
	ErrMissingGroupIDAndEmail = errors.BadRequest("GROUP_ID_OR_EMAIL_REQUIRED", "An email address or group id is needed to list invites")
	ErrInviteNotFound         = errors.NotFound("NOT_FOUND", "Invite not found")

	emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
)

type Invite struct {
	ID      string
	Email   string
	GroupID string
	Code    string
}

func (i *Invite) Serialize() *pb.Invite {
	return &pb.Invite{
		Id:      i.ID,
		Email:   i.Email,
		GroupId: i.GroupID,
		Code:    i.Code,
	}
}

func (i *Invite) Key(ctx context.Context) string {
	key := fmt.Sprintf("invite:%s:%s", i.ID, i.Code)

	t, ok := tenant.FromContext(ctx)
	if !ok {
		return key
	}

	return fmt.Sprintf("%s/%s", t, key)
}

func (i *Invite) Index(ctx context.Context) string {
	key := fmt.Sprintf("group:%s:%s", i.GroupID, i.Email)

	t, ok := tenant.FromContext(ctx)
	if !ok {
		return key
	}

	return fmt.Sprintf("%s/%s", t, key)
}

func (i *Invite) Marshal() []byte {
	b, _ := json.Marshal(i)
	return b
}

func (i *Invite) Unmarshal(b []byte) error {
	return json.Unmarshal(b, &i)
}

type Invites struct{}

// schema
// Read: id/code
// List: group/email

// Create an invite
func (i *Invites) Create(ctx context.Context, req *pb.CreateRequest, rsp *pb.CreateResponse) error {
	_, ok := auth.AccountFromContext(ctx)
	if !ok {
		errors.Unauthorized("UNAUTHORIZED", "Unauthorized")
	}
	// validate the request
	if len(req.GroupId) == 0 {
		return ErrMissingGroupID
	}
	if len(req.Email) == 0 {
		return ErrMissingEmail
	}
	if !isEmailValid(req.Email) {
		return ErrInvalidEmail
	}

	// construct the invite and write to the db
	invite := &Invite{
		ID:      uuid.New().String(),
		Code:    generateCode(),
		GroupID: req.GroupId,
		Email:   strings.ToLower(req.Email),
	}

	// id/val
	key := invite.Key(ctx)
	// get group key
	gkey := invite.Index(ctx)

	// TODO: Use the micro/micro/v3/service/sync interface to lock

	// write the first record
	if err := store.Write(store.NewRecord(key, invite)); err != nil {
		logger.Errorf("Error writing to the store: %v", err)
		return errors.InternalServerError("DATABASE_ERROR", "Error connecting to the database")
	}

	// write the group record
	if err := store.Write(store.NewRecord(gkey, invite)); err != nil {
		logger.Errorf("Error writing to the store: %v", err)
		return errors.InternalServerError("DATABASE_ERROR", "Error connecting to the database")
	}

	// serialize the response
	rsp.Invite = invite.Serialize()
	return nil
}

// Read an invite using ID or code
func (i *Invites) Read(ctx context.Context, req *pb.ReadRequest, rsp *pb.ReadResponse) error {
	_, ok := auth.AccountFromContext(ctx)
	if !ok {
		errors.Unauthorized("UNAUTHORIZED", "Unauthorized")
	}
	// validate the request
	if len(req.Id) == 0 && len(req.Code) == 0 {
		return ErrMissingIDAndCode
	}

	var recs []*store.Record
	var err error

	// create a pseudo invite
	invite := &Invite{
		ID:   req.Id,
		Code: req.Code,
	}

	if len(req.Id) > 0 && len(req.Code) > 0 {
		recs, err = store.Read(invite.Key(ctx), store.ReadLimit(1))
	} else if len(req.Id) > 0 {
		recs, err = store.Read(invite.Key(ctx), store.ReadLimit(1), store.ReadPrefix())
	} else if len(req.Code) > 0 {
		// create a code suffix key
		key := ":" + req.Code
		// read all the keys with the given code
		// TODO: potential race where if the code is not random
		// we read it for the wrong user e.g if two tenants generate the same code
		r, lerr := store.Read(key, store.ReadLimit(1), store.ReadSuffix())

		// now scan for the prefix
		prefix := "invite:"

		// additional prefix for the tenant
		if t, ok := tenant.FromContext(ctx); ok {
			prefix = t + "/" + prefix
		}

		// scan for the key we're looking for
		for _, rec := range r {
			// skip the missing prefix
			if !strings.HasPrefix(rec.Key, prefix) {
				continue
			}

			// skip missing suffix
			if !strings.HasSuffix(rec.Key, key) {
				continue
			}

			// save the record
			recs = append(recs, rec)
			break
		}

		// set the error
		// TODO: maybe just process this
		err = lerr
	}

	// check if there are any records
	if err == store.ErrNotFound || len(recs) == 0 {
		return ErrInviteNotFound
	}

	// check the error
	if err != nil {
		logger.Errorf("Error reading from the store: %v", err)
		return errors.InternalServerError("DATABASE_ERROR", "Error connecting to the database")
	}

	// unmarshal the invite
	invite.Unmarshal(recs[0].Value)

	// serialize the response
	rsp.Invite = invite.Serialize()
	return nil
}

// List invited for a group or specific email
func (i *Invites) List(ctx context.Context, req *pb.ListRequest, rsp *pb.ListResponse) error {
	_, ok := auth.AccountFromContext(ctx)
	if !ok {
		errors.Unauthorized("UNAUTHORIZED", "Unauthorized")
	}
	// validate the request
	if len(req.Email) == 0 && len(req.GroupId) == 0 {
		return ErrMissingGroupIDAndEmail
	}

	invite := &Invite{
		GroupID: req.GroupId,
		Email:   req.Email,
	}

	var recs []*store.Record
	var err error

	if len(invite.GroupID) > 0 && len(invite.Email) > 0 {
		key := invite.Index(ctx)
		recs, err = store.Read(key, store.ReadLimit(1))
	} else if len(invite.GroupID) > 0 {
		key := invite.Index(ctx)
		recs, err = store.Read(key, store.ReadPrefix())
	} else if len(invite.Email) > 0 {
		// create a email suffix key
		key := ":" + invite.Email
		// read all the keys with the given code
		r, lerr := store.Read(key, store.ReadSuffix())

		// now scan for the prefix
		prefix := "group:"

		// additional prefix for the tenant
		if t, ok := tenant.FromContext(ctx); ok {
			prefix = t + "/" + prefix
		}

		// scan for the key we're looking for
		for _, rec := range r {
			// skip the missing prefix
			if !strings.HasPrefix(rec.Key, prefix) {
				continue
			}

			// skip missing suffix
			if !strings.HasSuffix(rec.Key, key) {
				continue
			}

			// save the record
			recs = append(recs, rec)
		}

		// set the error
		// TODO: maybe just process this
		err = lerr
	}

	// no records found
	if err == store.ErrNotFound {
		return nil
	}

	// check the error
	if err != nil {
		logger.Errorf("Error reading from the store: %v", err)
		return errors.InternalServerError("DATABASE_ERROR", "Error connecting to the database")
	}

	// return response
	for _, rec := range recs {
		invite := new(Invite)
		invite.Unmarshal(rec.Value)
		rsp.Invites = append(rsp.Invites, invite.Serialize())
	}

	return nil
}

// Delete an invite
func (i *Invites) Delete(ctx context.Context, req *pb.DeleteRequest, rsp *pb.DeleteResponse) error {
	_, ok := auth.AccountFromContext(ctx)
	if !ok {
		errors.Unauthorized("UNAUTHORIZED", "Unauthorized")
	}
	// validate the request
	if len(req.Id) == 0 {
		return ErrMissingID
	}

	invite := &Invite{ID: req.Id}
	key := invite.Key(ctx)

	// check for the existing invite value
	recs, err := store.Read(key, store.ReadLimit(1), store.ReadPrefix())
	if err == store.ErrNotFound || len(recs) == 0 {
		return nil
	} else if err != nil {
		logger.Errorf("Error connecting to DB: %v", err)
		return errors.InternalServerError("DB_ERROR", "Error connecting to DB")
	}

	// unmarshal the existing invite
	invite.Unmarshal(recs[0].Value)
	if invite.ID != req.Id {
		return nil
	}

	// delete the record by id
	store.Delete(invite.Key(ctx))

	// delete the record by group id
	store.Delete(invite.Index(ctx))

	return nil
}

// isEmailValid checks if the email provided passes the required structure and length.
func isEmailValid(e string) bool {
	if len(e) < 3 && len(e) > 254 {
		return false
	}
	return emailRegex.MatchString(e)
}

// generateCode generates a random 8 digit code
func generateCode() string {
	v := rand.Intn(89999999) + 10000000
	return strconv.Itoa(v)
}
