package handler

import (
	"context"
	"time"

	"github.com/micro/micro/v3/service/auth"
	gorm2 "github.com/micro/services/pkg/gorm"
	"gorm.io/gorm"

	"github.com/google/uuid"
	"github.com/micro/micro/v3/service/errors"
	"github.com/micro/micro/v3/service/logger"
	pb "github.com/micro/services/seen/proto"
)

var (
	ErrMissingUserID       = errors.BadRequest("MISSING_USER_ID", "Missing UserID")
	ErrMissingResourceID   = errors.BadRequest("MISSING_RESOURCE_ID", "Missing ResourceID")
	ErrMissingResourceIDs  = errors.BadRequest("MISSING_RESOURCE_IDS", "Missing ResourceIDs")
	ErrMissingResourceType = errors.BadRequest("MISSING_RESOURCE_TYPE", "Missing ResourceType")
	ErrStore               = errors.InternalServerError("STORE_ERROR", "Error connecting to the store")
)

type Seen struct {
	gorm2.Helper
}

type SeenInstance struct {
	ID           string
	UserID       string `gorm:"uniqueIndex:user_resource"`
	ResourceID   string `gorm:"uniqueIndex:user_resource"`
	ResourceType string `gorm:"uniqueIndex:user_resource"`
	Timestamp    time.Time
}

// Set a resource as seen by a user. If no timestamp is provided, the current time is used.
func (s *Seen) Set(ctx context.Context, req *pb.SetRequest, rsp *pb.SetResponse) error {
	_, ok := auth.AccountFromContext(ctx)
	if !ok {
		errors.Unauthorized("UNAUTHORIZED", "Unauthorized")
	}
	// validate the request
	if len(req.UserId) == 0 {
		return ErrMissingUserID
	}
	if len(req.ResourceId) == 0 {
		return ErrMissingResourceID
	}
	if len(req.ResourceType) == 0 {
		return ErrMissingResourceType
	}

	// default the timestamp
	if req.Timestamp == 0 {
		req.Timestamp = time.Now().Unix()
	}

	// find the resource
	instance := SeenInstance{
		UserID:       req.UserId,
		ResourceID:   req.ResourceId,
		ResourceType: req.ResourceType,
	}
	db, err := s.GetDBConn(ctx)
	if err != nil {
		logger.Errorf("Error connecting to DB: %v", err)
		return errors.InternalServerError("DB_ERROR", "Error connecting to DB")
	}
	if err := db.Where(&instance).First(&instance).Error; err == gorm.ErrRecordNotFound {
		instance.ID = uuid.New().String()
	} else if err != nil {
		logger.Errorf("Error with store: %v", err)
		return ErrStore
	}

	// update the resource
	instance.Timestamp = time.Unix(req.Timestamp, 0)
	if err := db.Save(&instance).Error; err != nil {
		logger.Errorf("Error with store: %v", err)
		return ErrStore
	}

	return nil
}

// Unset a resource as seen, used in cases where a user viewed a resource but wants to override
// this so they remember to action it in the future, e.g. "Mark this as unread".
func (s *Seen) Unset(ctx context.Context, req *pb.UnsetRequest, rsp *pb.UnsetResponse) error {
	_, ok := auth.AccountFromContext(ctx)
	if !ok {
		errors.Unauthorized("UNAUTHORIZED", "Unauthorized")
	}
	// validate the request
	if len(req.UserId) == 0 {
		return ErrMissingUserID
	}
	if len(req.ResourceId) == 0 {
		return ErrMissingResourceID
	}
	if len(req.ResourceType) == 0 {
		return ErrMissingResourceType
	}

	db, err := s.GetDBConn(ctx)
	if err != nil {
		logger.Errorf("Error connecting to DB: %v", err)
		return errors.InternalServerError("DB_ERROR", "Error connecting to DB")
	}
	// delete the object from the store
	err = db.Delete(SeenInstance{}, SeenInstance{
		UserID:       req.UserId,
		ResourceID:   req.ResourceId,
		ResourceType: req.ResourceType,
	}).Error
	if err != nil {
		logger.Errorf("Error with store: %v", err)
		return ErrStore
	}

	return nil
}

// Read returns the timestamps at which various resources were seen by a user. If no timestamp
// is returned for a given resource_id, it indicates that resource has not yet been seen by the
// user.
func (s *Seen) Read(ctx context.Context, req *pb.ReadRequest, rsp *pb.ReadResponse) error {
	_, ok := auth.AccountFromContext(ctx)
	if !ok {
		errors.Unauthorized("UNAUTHORIZED", "Unauthorized")
	}
	// validate the request
	if len(req.UserId) == 0 {
		return ErrMissingUserID
	}
	if len(req.ResourceIds) == 0 {
		return ErrMissingResourceIDs
	}
	if len(req.ResourceType) == 0 {
		return ErrMissingResourceType
	}

	db, err := s.GetDBConn(ctx)
	if err != nil {
		logger.Errorf("Error connecting to DB: %v", err)
		return errors.InternalServerError("DB_ERROR", "Error connecting to DB")
	}
	// query the store
	q := db.Where(SeenInstance{UserID: req.UserId, ResourceType: req.ResourceType})
	q = q.Where("resource_id IN (?)", req.ResourceIds)
	var data []SeenInstance
	if err := q.Find(&data).Error; err != nil {
		logger.Errorf("Error with store: %v", err)
		return ErrStore
	}

	// serialize the response
	rsp.Timestamps = make(map[string]int64, len(data))
	for _, i := range data {
		rsp.Timestamps[i.ResourceID] = i.Timestamp.Unix()
	}

	return nil
}
