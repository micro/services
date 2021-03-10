package handler_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/micro/services/seen/handler"
	pb "github.com/micro/services/seen/proto"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func testHandler(t *testing.T) *handler.Seen {
	// connect to the database
	addr := os.Getenv("POSTGRES_URL")
	if len(addr) == 0 {
		addr = "postgresql://postgres@localhost:5432/postgres?sslmode=disable"
	}
	db, err := gorm.Open(postgres.Open(addr), &gorm.Config{})
	if err != nil {
		t.Fatalf("Error connecting to database: %v", err)
	}

	// migrate the database
	if err := db.AutoMigrate(&handler.SeenInstance{}); err != nil {
		t.Fatalf("Error migrating database: %v", err)
	}

	// clean any data from a previous run
	if err := db.Exec("TRUNCATE TABLE seen_instances CASCADE").Error; err != nil {
		t.Fatalf("Error cleaning database: %v", err)
	}

	return &handler.Seen{DB: db}
}

func TestSet(t *testing.T) {
	tt := []struct {
		Name         string
		UserID       string
		ResourceType string
		ResourceID   string
		Timestamp    *timestamppb.Timestamp
		Error        error
	}{
		{
			Name:         "MissingUserID",
			ResourceType: "message",
			ResourceID:   uuid.New().String(),
			Error:        handler.ErrMissingUserID,
		},
		{
			Name:         "MissingResourceID",
			UserID:       uuid.New().String(),
			ResourceType: "message",
			Error:        handler.ErrMissingResourceID,
		},
		{
			Name:       "MissingResourceType",
			UserID:     uuid.New().String(),
			ResourceID: uuid.New().String(),
			Error:      handler.ErrMissingResourceType,
		},
		{
			Name:         "WithTimetamp",
			UserID:       uuid.New().String(),
			ResourceID:   uuid.New().String(),
			ResourceType: "message",
			Timestamp:    timestamppb.New(time.Now().Add(time.Minute * -5)),
		},
		{
			Name:         "WithoutTimetamp",
			UserID:       uuid.New().String(),
			ResourceID:   uuid.New().String(),
			ResourceType: "message",
		},
		{
			Name:         "WithUpdatedTimetamp",
			UserID:       uuid.New().String(),
			ResourceID:   uuid.New().String(),
			ResourceType: "message",
			Timestamp:    timestamppb.New(time.Now().Add(time.Minute * -3)),
		},
	}

	h := testHandler(t)
	for _, tc := range tt {
		t.Run(tc.Name, func(t *testing.T) {
			err := h.Set(context.TODO(), &pb.SetRequest{
				UserId:       tc.UserID,
				ResourceId:   tc.ResourceID,
				ResourceType: tc.ResourceType,
				Timestamp:    tc.Timestamp,
			}, &pb.SetResponse{})

			assert.Equal(t, tc.Error, err)
		})
	}
}
func TestUnset(t *testing.T) {
	// seed some test data
	h := testHandler(t)
	seed := &pb.SetRequest{
		UserId:       uuid.New().String(),
		ResourceId:   uuid.New().String(),
		ResourceType: "message",
	}
	err := h.Set(context.TODO(), seed, &pb.SetResponse{})
	assert.NoError(t, err)

	tt := []struct {
		Name         string
		UserID       string
		ResourceType string
		ResourceID   string
		Error        error
	}{
		{
			Name:         "MissingUserID",
			ResourceType: "message",
			ResourceID:   uuid.New().String(),
			Error:        handler.ErrMissingUserID,
		},
		{
			Name:         "MissingResourceID",
			UserID:       uuid.New().String(),
			ResourceType: "message",
			Error:        handler.ErrMissingResourceID,
		},
		{
			Name:       "MissingResourceType",
			UserID:     uuid.New().String(),
			ResourceID: uuid.New().String(),
			Error:      handler.ErrMissingResourceType,
		},
		{
			Name:         "Exists",
			UserID:       seed.UserId,
			ResourceID:   seed.ResourceId,
			ResourceType: seed.ResourceType,
		},
		{
			Name:         "Repeat",
			UserID:       seed.UserId,
			ResourceID:   seed.ResourceId,
			ResourceType: seed.ResourceType,
		},
	}

	for _, tc := range tt {
		t.Run(tc.Name, func(t *testing.T) {
			err := h.Unset(context.TODO(), &pb.UnsetRequest{
				UserId:       tc.UserID,
				ResourceId:   tc.ResourceID,
				ResourceType: tc.ResourceType,
			}, &pb.UnsetResponse{})

			assert.Equal(t, tc.Error, err)
		})
	}
}

func TestRead(t *testing.T) {
	tn := time.Now()
	h := testHandler(t)

	// seed some test data
	td := []struct {
		UserID       string
		ResourceID   string
		ResourceType string
		Timestamp    *timestamppb.Timestamp
	}{
		{
			UserID:       "user-1",
			ResourceID:   "message-1",
			ResourceType: "message",
			Timestamp:    timestamppb.New(tn.Add(time.Minute * -10)),
		},
		{
			UserID:       "user-1",
			ResourceID:   "message-1",
			ResourceType: "message",
			Timestamp:    timestamppb.New(tn),
		},
		{
			UserID:       "user-1",
			ResourceID:   "message-2",
			ResourceType: "message",
			Timestamp:    timestamppb.New(tn.Add(time.Minute * -10)),
		},
		{
			UserID:       "user-1",
			ResourceID:   "notification-1",
			ResourceType: "notification",
			Timestamp:    timestamppb.New(tn.Add(time.Minute * -10)),
		},
		{
			UserID:       "user-2",
			ResourceID:   "message-3",
			ResourceType: "message",
			Timestamp:    timestamppb.New(tn.Add(time.Minute * -10)),
		},
	}
	for _, d := range td {
		assert.NoError(t, h.Set(context.TODO(), &pb.SetRequest{
			UserId:       d.UserID,
			ResourceId:   d.ResourceID,
			ResourceType: d.ResourceType,
			Timestamp:    d.Timestamp,
		}, &pb.SetResponse{}))
	}

	// check only the requested values are returned
	var rsp pb.ReadResponse
	err := h.Read(context.TODO(), &pb.ReadRequest{
		UserId:       "user-1",
		ResourceType: "message",
		ResourceIds:  []string{"message-1", "message-2", "message-3"},
	}, &rsp)
	assert.NoError(t, err)
	assert.Len(t, rsp.Timestamps, 2)

	if v := rsp.Timestamps["message-1"]; v != nil {
		assert.True(t, v.AsTime().Equal(tn))
	} else {
		t.Errorf("Expected a timestamp for message-1")
	}

	if v := rsp.Timestamps["message-2"]; v != nil {
		assert.True(t, v.AsTime().Equal(tn.Add(time.Minute*-10)))
	} else {
		t.Errorf("Expected a timestamp for message-2")
	}

	// unsetting a resource should remove it from the list
	err = h.Unset(context.TODO(), &pb.UnsetRequest{
		UserId:       "user-1",
		ResourceId:   "message-2",
		ResourceType: "message",
	}, &pb.UnsetResponse{})
	assert.NoError(t, err)

	rsp = pb.ReadResponse{}
	err = h.Read(context.TODO(), &pb.ReadRequest{
		UserId:       "user-1",
		ResourceType: "message",
		ResourceIds:  []string{"message-1", "message-2", "message-3"},
	}, &rsp)
	assert.NoError(t, err)
	assert.Len(t, rsp.Timestamps, 1)
}
