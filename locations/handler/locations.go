package handler

import (
	"context"
	"sync"
	"time"

	"github.com/google/uuid"
	geo "github.com/hailocab/go-geoindex"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
	"gorm.io/gorm"

	"github.com/micro/micro/v3/service/errors"
	"github.com/micro/micro/v3/service/logger"
	"github.com/micro/services/locations/model"
	pb "github.com/micro/services/locations/proto"
)

var (
	ErrMissingLocations = errors.BadRequest("MISSING_LOCATIONS", "One or more locations are required")
	ErrMissingLatitude  = errors.BadRequest("MISSING_LATITUDE", "Latitude is required")
	ErrMissingLongitude = errors.BadRequest("MISSING_LONGITUDE", "Longitude is required")
	ErrMissingUserID    = errors.BadRequest("MISSING_USER_ID", "UserID is required")
	ErrMissingUserIDs   = errors.BadRequest("MISSING_USER_IDS", "One or more UserIDs are required")
	ErrMissingBefore    = errors.BadRequest("MISSING_BEFORE", "Before timestamp is required")
	ErrMissingAfter     = errors.BadRequest("MISSING_AFTER", "After timestamp is required")
	ErrMissingRadius    = errors.BadRequest("MISSING_RADIUS", "Radius is required")
)

type Locations struct {
	sync.RWMutex

	Geoindex *geo.PointsIndex
	DB       *gorm.DB
}

// Save a set of locations
func (l *Locations) Save(ctx context.Context, req *pb.SaveRequest, rsp *pb.SaveResponse) error {
	// validate the request
	if len(req.Locations) == 0 {
		return ErrMissingLocations
	}
	for _, l := range req.Locations {
		if l.Latitude == nil {
			return ErrMissingLatitude
		}
		if l.Longitude == nil {
			return ErrMissingLongitude
		}
		if len(l.UserId) == 0 {
			return ErrMissingUserID
		}
	}

	// construct the database objects
	ls := make([]*model.Location, len(req.Locations))
	for i, lc := range req.Locations {
		loc := &model.Location{
			ID:        uuid.New().String(),
			UserID:    lc.UserId,
			Latitude:  lc.Latitude.Value,
			Longitude: lc.Longitude.Value,
		}
		if lc.Timestamp != nil {
			loc.Timestamp = lc.Timestamp.AsTime()
		} else {
			loc.Timestamp = time.Now()
		}
		ls[i] = loc
	}

	// write to the database
	if err := l.DB.Create(ls).Error; err != nil {
		logger.Errorf("Error writing to the database: %v", err)
		return errors.InternalServerError("DATABASE_ERROR", "Error writing to the database")
	}

	// write to the geoindex
	l.Lock()
	defer l.Unlock()
	for _, lc := range ls {
		l.Geoindex.Add(lc)
	}
	return nil
}

// Last locations for a set of users
func (l *Locations) Last(ctx context.Context, req *pb.LastRequest, rsp *pb.ListResponse) error {
	// validate the request
	if req.UserIds == nil {
		return ErrMissingUserIDs
	}

	// query the database
	q := l.DB.Raw("SELECT DISTINCT ON (user_id) user_id, timestamp, latitude, longitude FROM locations WHERE user_id IN (?) ORDER BY user_id, timestamp DESC", req.UserIds)
	var locs []*model.Location
	if err := q.Find(&locs).Error; err != nil {
		logger.Errorf("Error reading from the database: %v", err)
		return errors.InternalServerError("DATABASE_ERROR", "Error reading from the database")
	}

	// serialize the result
	rsp.Locations = serializeLocations(locs)
	return nil
}

// Near returns the locations near a point
func (l *Locations) Near(ctx context.Context, req *pb.NearRequest, rsp *pb.ListResponse) error {
	// validate the request
	if req.Latitude == nil {
		return ErrMissingLatitude
	}
	if req.Longitude == nil {
		return ErrMissingLongitude
	}
	if req.Radius == nil {
		return ErrMissingRadius
	}

	// query the geoindex
	l.RLock()
	p := geo.NewGeoPoint("query", req.Latitude.Value, req.Longitude.Value)
	result := l.Geoindex.PointsWithin(p, geo.Km(req.Radius.Value), func(p geo.Point) bool {
		return true
	})
	l.RUnlock()

	// serialize the result
	locs := make([]*model.Location, len(result))
	for i, r := range result {
		locs[i] = r.(*model.Location)
	}
	rsp.Locations = serializeLocations(locs)
	return nil
}

// Read locations for a group of users between two points in time
func (l *Locations) Read(ctx context.Context, req *pb.ReadRequest, rsp *pb.ListResponse) error {
	// validate the request
	if len(req.UserIds) == 0 {
		return ErrMissingUserIDs
	}
	if req.Before == nil {
		return ErrMissingBefore
	}
	if req.After == nil {
		return ErrMissingAfter
	}

	// construct the request
	q := l.DB.Model(&model.Location{})
	q = q.Order("timestamp ASC")
	q = q.Where("user_id IN (?) AND timestamp > ? AND timestamp < ?", req.UserIds, req.After.AsTime(), req.Before.AsTime())
	var locs []*model.Location
	if err := q.Find(&locs).Error; err != nil {
		logger.Errorf("Error reading from the database: %v", err)
		return errors.InternalServerError("DATABASE_ERROR", "Error reading from the database")
	}

	// serialize the result
	rsp.Locations = serializeLocations(locs)
	return nil
}

func serializeLocations(locs []*model.Location) []*pb.Location {
	rsp := make([]*pb.Location, len(locs))
	for i, l := range locs {
		rsp[i] = &pb.Location{
			UserId:    l.UserID,
			Latitude:  &wrapperspb.DoubleValue{Value: l.Latitude},
			Longitude: &wrapperspb.DoubleValue{Value: l.Longitude},
			Timestamp: timestamppb.New(l.Timestamp),
		}
	}
	return rsp
}
