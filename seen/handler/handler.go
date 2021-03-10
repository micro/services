package handler

import (
	"context"
	"time"

	"gorm.io/gorm"

	"github.com/google/uuid"
	"github.com/micro/micro/v3/service/errors"
	"github.com/micro/micro/v3/service/logger"
	pb "github.com/micro/services/seen/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var (
	ErrMissingUserID       = errors.BadRequest("MISSING_USER_ID", "Missing UserID")
	ErrMissingResourceID   = errors.BadRequest("MISSING_RESOURCE_ID", "Missing ResourceID")
	ErrMissingResourceIDs  = errors.BadRequest("MISSING_RESOURCE_IDS", "Missing ResourceIDs")
	ErrMissingResourceType = errors.BadRequest("MISSING_RESOURCE_TYPE", "Missing ResourceType")
	ErrStore               = errors.InternalServerError("STORE_ERROR", "Error connecting to the store")
)

type Seen struct {
	DB *gorm.DB
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
	if req.Timestamp == nil {
		req.Timestamp = timestamppb.New(time.Now())
	}

	// find the resource
	instance := SeenInstance{
		UserID:       req.UserId,
		ResourceID:   req.ResourceId,
		ResourceType: req.ResourceType,
	}
	if err := s.DB.Where(&instance).First(&instance).Error; err == gorm.ErrRecordNotFound {
		instance.ID = uuid.New().String()
	} else if err != nil {
		logger.Errorf("Error with store: %v", err)
		return ErrStore
	}

	// update the resource
	instance.Timestamp = req.Timestamp.AsTime()
	if err := s.DB.Save(&instance).Error; err != nil {
		logger.Errorf("Error with store: %v", err)
		return ErrStore
	}

	return nil
}

// Unset a resource as seen, used in cases where a user viewed a resource but wants to override
// this so they remember to action it in the future, e.g. "Mark this as unread".
func (s *Seen) Unset(ctx context.Context, req *pb.UnsetRequest, rsp *pb.UnsetResponse) error {
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

	// delete the object from the store
	err := s.DB.Delete(SeenInstance{}, SeenInstance{
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

	// query the store
	q := s.DB.Where(SeenInstance{UserID: req.UserId, ResourceType: req.ResourceType})
	q = q.Where("resource_id IN (?)", req.ResourceIds)
	var data []SeenInstance
	if err := q.Find(&data).Error; err != nil {
		logger.Errorf("Error with store: %v", err)
		return ErrStore
	}

	// serialize the response
	rsp.Timestamps = make(map[string]*timestamppb.Timestamp, len(data))
	for _, i := range data {
		rsp.Timestamps[i.ResourceID] = timestamppb.New(i.Timestamp)
	}

	return nil
}
