package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/micro/micro/v3/service/auth"
	"github.com/micro/micro/v3/service/errors"
	"github.com/micro/micro/v3/service/logger"
	"github.com/micro/micro/v3/service/store"
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

type Seen struct{}

type Record struct {
	ID           string
	UserID       string
	ResourceID   string
	ResourceType string
	Timestamp    time.Time
}

func (r *Record) Key() string {
	return fmt.Sprintf("%s:%s:%s", r.UserID, r.ResourceType, r.ResourceID)
}

func (r *Record) Marshal() []byte {
	b, _ := json.Marshal(r)
	return b
}

func (r *Record) Unmarshal(b []byte) error {
	return json.Unmarshal(b, &r)
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
	if req.Timestamp == nil {
		req.Timestamp = timestamppb.New(time.Now())
	}

	// find the resource
	instance := &Record{
		UserID:       req.UserId,
		ResourceID:   req.ResourceId,
		ResourceType: req.ResourceType,
	}

	_, err := store.Read(instance.Key(), store.ReadLimit(1))
	if err == store.ErrNotFound {
		instance.ID = uuid.New().String()
	} else if err != nil {
		logger.Errorf("Error with store: %v", err)
		return ErrStore
	}

	// update the resource
	instance.Timestamp = req.Timestamp.AsTime()

	if err := store.Write(&store.Record{
		Key:   instance.Key(),
		Value: instance.Marshal(),
	}); err != nil {
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

	instance := &Record{
		UserID:       req.UserId,
		ResourceID:   req.ResourceId,
		ResourceType: req.ResourceType,
	}

	// delete the object from the store
	if err := store.Delete(instance.Key()); err != nil {
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

	// create a key prefix
	key := fmt.Sprintf("%s:%s:", req.UserId, req.ResourceType)

	var recs []*store.Record
	var err error

	// get the records for the resource type
	if len(req.ResourceIds) == 1 {
		// read the key itself
		key = key + req.ResourceIds[0]
		recs, err = store.Read(key, store.ReadLimit(1))
	} else {
		// otherwise read the prefix
		recs, err = store.Read(key, store.ReadPrefix())
	}

	if err != nil {
		logger.Errorf("Error with store: %v", err)
		return ErrStore
	}

	// make an id map
	ids := make(map[string]bool)

	for _, id := range req.ResourceIds {
		ids[id] = true
	}

	// make the map
	rsp.Timestamps = make(map[string]*timestamppb.Timestamp)

	// range over records for the user/resource type
	// TODO: add some sort of filter query in store
	for _, rec := range recs {
		// get id
		parts := strings.Split(rec.Key, ":")
		id := parts[2]

		fmt.Println("checking record", rec.Key, id)

		if ok := ids[id]; !ok {
			continue
		}

		// add the timestamp for the record
		r := new(Record)
		r.Unmarshal(rec.Value)
		rsp.Timestamps[id] = timestamppb.New(r.Timestamp)
	}

	return nil
}
