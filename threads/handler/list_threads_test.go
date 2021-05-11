package handler_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/micro/services/threads/handler"
	pb "github.com/micro/services/threads/proto"
	"github.com/stretchr/testify/assert"
)

func TestListThreads(t *testing.T) {
	h := testHandler(t)

	// seed some data
	var cRsp1 pb.CreateThreadResponse
	err := h.CreateThread(microAccountCtx(), &pb.CreateThreadRequest{
		Topic: "HelloWorld", GroupId: uuid.New().String(),
	}, &cRsp1)
	if err != nil {
		t.Fatalf("Error creating thread: %v", err)
		return
	}
	var cRsp2 pb.CreateThreadResponse
	err = h.CreateThread(microAccountCtx(), &pb.CreateThreadRequest{
		Topic: "FooBar", GroupId: uuid.New().String(),
	}, &cRsp2)
	if err != nil {
		t.Fatalf("Error creating thread: %v", err)
		return
	}

	t.Run("MissingGroupID", func(t *testing.T) {
		var rsp pb.ListThreadsResponse
		err := h.ListThreads(microAccountCtx(), &pb.ListThreadsRequest{}, &rsp)
		assert.Equal(t, handler.ErrMissingGroupID, err)
		assert.Nil(t, rsp.Threads)
	})

	t.Run("Valid", func(t *testing.T) {
		var rsp pb.ListThreadsResponse
		err := h.ListThreads(microAccountCtx(), &pb.ListThreadsRequest{
			GroupId: cRsp1.Thread.GroupId,
		}, &rsp)

		assert.NoError(t, err)
		if len(rsp.Threads) != 1 {
			t.Fatalf("Expected 1 thread to be returned, got %v", len(rsp.Threads))
			return
		}

		assertThreadsMatch(t, cRsp1.Thread, rsp.Threads[0])
	})
}
