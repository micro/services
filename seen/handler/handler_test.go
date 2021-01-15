package handler_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/google/uuid"
	"github.com/micro/services/seen/domain"
	"github.com/micro/services/seen/handler"
	pb "github.com/micro/services/seen/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func newHandler() *handler.Seen {
	return &handler.Seen{
		Domain: domain.New(),
	}
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
	}

	h := newHandler()
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
	h := newHandler()
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
	h := newHandler()

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
