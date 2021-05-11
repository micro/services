package handler_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/micro/services/threads/handler"
	pb "github.com/micro/services/threads/proto"
	"github.com/stretchr/testify/assert"
)

func TestReadThread(t *testing.T) {
	h := testHandler(t)

	// seed some data
	var cRsp pb.CreateThreadResponse
	err := h.CreateThread(microAccountCtx(), &pb.CreateThreadRequest{
		Topic: "HelloWorld", GroupId: uuid.New().String(),
	}, &cRsp)
	if err != nil {
		t.Fatalf("Error creating thread: %v", err)
		return
	}

	tt := []struct {
		Name    string
		ID      string
		GroupID string
		Error   error
		Result  *pb.Thread
	}{
		{
			Name:  "MissingID",
			Error: handler.ErrMissingID,
		},
		{
			Name:  "IncorrectID",
			ID:    uuid.New().String(),
			Error: handler.ErrNotFound,
		},
		{
			Name:   "FoundUsingIDOnly",
			ID:     cRsp.Thread.Id,
			Result: cRsp.Thread,
		},
		{
			Name:    "IncorrectGroupID",
			ID:      cRsp.Thread.Id,
			Error:   handler.ErrNotFound,
			GroupID: uuid.New().String(),
		},
	}
	for _, tc := range tt {
		t.Run(tc.Name, func(t *testing.T) {
			var rsp pb.ReadThreadResponse
			err := h.ReadThread(microAccountCtx(), &pb.ReadThreadRequest{
				Id: tc.ID, GroupId: tc.GroupID,
			}, &rsp)
			assert.Equal(t, tc.Error, err)

			if tc.Result == nil {
				assert.Nil(t, rsp.Thread)
			} else {
				assertThreadsMatch(t, tc.Result, rsp.Thread)
			}
		})
	}
}
