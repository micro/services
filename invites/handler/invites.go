package handler

import (
	"context"
	"math/rand"
	"regexp"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/micro/micro/v3/service/errors"
	"github.com/micro/micro/v3/service/logger"
	pb "github.com/micro/services/invites/proto"
	"gorm.io/gorm"
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
	Email   string `gorm:"uniqueIndex:group_email"`
	GroupID string `gorm:"uniqueIndex:group_email"`
	Code    string `gorm:"uniqueIndex"`
}

func (i *Invite) Serialize() *pb.Invite {
	return &pb.Invite{
		Id:      i.ID,
		Email:   i.Email,
		GroupId: i.GroupID,
		Code:    i.Code,
	}
}

type Invites struct {
	DB *gorm.DB
}

// Create an invite
func (i *Invites) Create(ctx context.Context, req *pb.CreateRequest, rsp *pb.CreateResponse) error {
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
		Email:   req.Email,
	}
	if err := i.DB.Create(invite).Error; err != nil && strings.Contains(err.Error(), "group_email") {
	} else if err != nil {
		logger.Errorf("Error writing to the store: %v", err)
		return errors.InternalServerError("DATABASE_ERROR", "Error connecting to the database")
	}

	// serialize the response
	rsp.Invite = invite.Serialize()
	return nil
}

// Read an invite using ID or code
func (i *Invites) Read(ctx context.Context, req *pb.ReadRequest, rsp *pb.ReadResponse) error {
	// validate the request
	var query Invite
	if req.Id != nil {
		query.ID = req.Id.Value
	} else if req.Code != nil {
		query.Code = req.Code.Value
	} else {
		return ErrMissingIDAndCode
	}

	// query the database
	var invite Invite
	if err := i.DB.Where(&query).First(&invite).Error; err == gorm.ErrRecordNotFound {
		return ErrInviteNotFound
	} else if err != nil {
		logger.Errorf("Error reading from the store: %v", err)
		return errors.InternalServerError("DATABASE_ERROR", "Error connecting to the database")
	}

	// serialize the response
	rsp.Invite = invite.Serialize()
	return nil
}

// List invited for a group or specific email
func (i *Invites) List(ctx context.Context, req *pb.ListRequest, rsp *pb.ListResponse) error {
	// validate the request
	if req.Email == nil && req.GroupId == nil {
		return ErrMissingGroupIDAndEmail
	}

	// construct the query
	var query Invite
	if req.GroupId != nil {
		query.GroupID = req.GroupId.Value
	}
	if req.Email != nil {
		query.Email = req.Email.Value
	}

	// query the database
	var invites []Invite
	if err := i.DB.Where(&query).Find(&invites).Error; err != nil {
		logger.Errorf("Error reading from the store: %v", err)
		return errors.InternalServerError("DATABASE_ERROR", "Error connecting to the database")
	}

	// serialize the response
	rsp.Invites = make([]*pb.Invite, len(invites))
	for i, inv := range invites {
		rsp.Invites[i] = inv.Serialize()
	}
	return nil
}

// Delete an invite
func (i *Invites) Delete(ctx context.Context, req *pb.DeleteRequest, rsp *pb.DeleteResponse) error {
	// validate the request
	if len(req.Id) == 0 {
		return ErrMissingID
	}

	// delete from the database
	if err := i.DB.Where(&Invite{ID: req.Id}).Delete(&Invite{}).Error; err != nil {
		logger.Errorf("Error deleting from the store: %v", err)
		return errors.InternalServerError("DATABASE_ERROR", "Error connecting to the database")
	}

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
