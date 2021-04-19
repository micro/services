package handler_test

import (
	"context"
	"database/sql"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/micro/micro/v3/service/auth"
	"github.com/micro/services/seen/handler"
	pb "github.com/micro/services/seen/proto"
	"github.com/stretchr/testify/assert"
)

func testHandler(t *testing.T) *handler.Seen {
	// connect to the database
	addr := os.Getenv("POSTGRES_URL")
	if len(addr) == 0 {
		addr = "postgresql://postgres@localhost:5432/postgres?sslmode=disable"
	}
	sqlDB, err := sql.Open("pgx", addr)
	if err != nil {
		t.Fatalf("Failed to open connection to DB %s", err)
	}

	// clean any data from a previous run
	if _, err := sqlDB.Exec("DROP TABLE IF EXISTS micro_seen_instances CASCADE"); err != nil {
		t.Fatalf("Error cleaning database: %v", err)
	}

	h := &handler.Seen{}
	h.DBConn(sqlDB).Migrations(&handler.SeenInstance{})
	return h
}

func TestSet(t *testing.T) {
	tt := []struct {
		Name         string
		UserID       string
		ResourceType string
		ResourceID   string
		Timestamp    int64
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
			Timestamp:    time.Now().Add(time.Minute * -5).Unix(),
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
			Timestamp:    time.Now().Add(time.Minute * -3).Unix(),
		},
	}

	h := testHandler(t)
	for _, tc := range tt {
		t.Run(tc.Name, func(t *testing.T) {
			err := h.Set(microAccountCtx(), &pb.SetRequest{
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
	err := h.Set(microAccountCtx(), seed, &pb.SetResponse{})
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
			err := h.Unset(microAccountCtx(), &pb.UnsetRequest{
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
		Timestamp    int64
	}{
		{
			UserID:       "user-1",
			ResourceID:   "message-1",
			ResourceType: "message",
			Timestamp:    tn.Add(time.Minute * -10).Unix(),
		},
		{
			UserID:       "user-1",
			ResourceID:   "message-1",
			ResourceType: "message",
			Timestamp:    tn.Unix(),
		},
		{
			UserID:       "user-1",
			ResourceID:   "message-2",
			ResourceType: "message",
			Timestamp:    tn.Add(time.Minute * -10).Unix(),
		},
		{
			UserID:       "user-1",
			ResourceID:   "notification-1",
			ResourceType: "notification",
			Timestamp:    tn.Add(time.Minute * -10).Unix(),
		},
		{
			UserID:       "user-2",
			ResourceID:   "message-3",
			ResourceType: "message",
			Timestamp:    tn.Add(time.Minute * -10).Unix(),
		},
	}
	for _, d := range td {
		assert.NoError(t, h.Set(microAccountCtx(), &pb.SetRequest{
			UserId:       d.UserID,
			ResourceId:   d.ResourceID,
			ResourceType: d.ResourceType,
			Timestamp:    d.Timestamp,
		}, &pb.SetResponse{}))
	}

	// check only the requested values are returned
	var rsp pb.ReadResponse
	err := h.Read(microAccountCtx(), &pb.ReadRequest{
		UserId:       "user-1",
		ResourceType: "message",
		ResourceIds:  []string{"message-1", "message-2", "message-3"},
	}, &rsp)
	assert.NoError(t, err)
	assert.Len(t, rsp.Timestamps, 2)

	if v := rsp.Timestamps["message-1"]; v != 0 {
		assert.Equal(t, microSecondTime(time.Unix(v, 0)), microSecondTime(tn))
	} else {
		t.Errorf("Expected a timestamp for message-1")
	}

	if v := rsp.Timestamps["message-2"]; v != 0 {
		assert.Equal(t, microSecondTime(time.Unix(v, 0)), microSecondTime(tn.Add(time.Minute*-10).UTC()))
	} else {
		t.Errorf("Expected a timestamp for message-2")
	}

	// unsetting a resource should remove it from the list
	err = h.Unset(microAccountCtx(), &pb.UnsetRequest{
		UserId:       "user-1",
		ResourceId:   "message-2",
		ResourceType: "message",
	}, &pb.UnsetResponse{})
	assert.NoError(t, err)

	rsp = pb.ReadResponse{}
	err = h.Read(microAccountCtx(), &pb.ReadRequest{
		UserId:       "user-1",
		ResourceType: "message",
		ResourceIds:  []string{"message-1", "message-2", "message-3"},
	}, &rsp)
	assert.NoError(t, err)
	assert.Len(t, rsp.Timestamps, 1)

}

// postgres has a resolution of 100microseconds so just test that it's accurate to the second
func microSecondTime(tt time.Time) time.Time {
	return time.Unix(tt.Unix(), int64(tt.Nanosecond()-tt.Nanosecond()%1000)).UTC()
}

func microAccountCtx() context.Context {
	return auth.ContextWithAccount(context.TODO(), &auth.Account{
		Issuer: "micro",
	})
}
