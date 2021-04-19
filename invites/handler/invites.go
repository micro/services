package handler

import (
	"context"
	"math/rand"
	"regexp"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/micro/micro/v3/service/auth"
	"github.com/micro/micro/v3/service/errors"
	"github.com/micro/micro/v3/service/logger"
	pb "github.com/micro/services/invites/proto"
	gorm2 "github.com/micro/services/pkg/gorm"
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
	gorm2.Helper
}

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
	db, err := i.GetDBConn(ctx)
	if err != nil {
		logger.Errorf("Error connecting to DB: %v", err)
		return errors.InternalServerError("DB_ERROR", "Error connecting to DB")
	}
	if err := db.Create(invite).Error; err != nil && !strings.Contains(err.Error(), "group_email") {
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
	var query Invite
	if req.Id != "" {
		query.ID = req.Id
	} else if req.Code != "" {
		query.Code = req.Code
	} else {
		return ErrMissingIDAndCode
	}

	db, err := i.GetDBConn(ctx)
	if err != nil {
		logger.Errorf("Error connecting to DB: %v", err)
		return errors.InternalServerError("DB_ERROR", "Error connecting to DB")
	}
	// query the database
	var invite Invite
	if err := db.Where(&query).First(&invite).Error; err == gorm.ErrRecordNotFound {
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
	_, ok := auth.AccountFromContext(ctx)
	if !ok {
		errors.Unauthorized("UNAUTHORIZED", "Unauthorized")
	}
	// validate the request
	if req.Email == "" && req.GroupId == "" {
		return ErrMissingGroupIDAndEmail
	}

	// construct the query
	var query Invite
	if req.GroupId != "" {
		query.GroupID = req.GroupId
	}
	if req.Email != "" {
		query.Email = strings.ToLower(req.Email)
	}

	db, err := i.GetDBConn(ctx)
	if err != nil {
		logger.Errorf("Error connecting to DB: %v", err)
		return errors.InternalServerError("DB_ERROR", "Error connecting to DB")
	}
	// query the database
	var invites []Invite
	if err := db.Where(&query).Find(&invites).Error; err != nil {
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
	_, ok := auth.AccountFromContext(ctx)
	if !ok {
		errors.Unauthorized("UNAUTHORIZED", "Unauthorized")
	}
	// validate the request
	if len(req.Id) == 0 {
		return ErrMissingID
	}

	db, err := i.GetDBConn(ctx)
	if err != nil {
		logger.Errorf("Error connecting to DB: %v", err)
		return errors.InternalServerError("DB_ERROR", "Error connecting to DB")
	}
	// delete from the database
	if err := db.Where(&Invite{ID: req.Id}).Delete(&Invite{}).Error; err != nil {
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
