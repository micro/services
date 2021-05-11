package handler_test

import (
	"testing"

	"github.com/micro/services/threads/handler"
	pb "github.com/micro/services/threads/proto"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestDeleteThread(t *testing.T) {
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

	t.Run("MissingID", func(t *testing.T) {
		err := h.DeleteThread(microAccountCtx(), &pb.DeleteThreadRequest{}, &pb.DeleteThreadResponse{})
		assert.Equal(t, handler.ErrMissingID, err)
	})

	t.Run("Valid", func(t *testing.T) {
		err := h.DeleteThread(microAccountCtx(), &pb.DeleteThreadRequest{
			Id: cRsp.Thread.Id,
		}, &pb.DeleteThreadResponse{})
		assert.NoError(t, err)

		err = h.ReadThread(microAccountCtx(), &pb.ReadThreadRequest{
			Id: cRsp.Thread.Id,
		}, &pb.ReadThreadResponse{})
		assert.Equal(t, handler.ErrNotFound, err)
	})

	t.Run("Retry", func(t *testing.T) {
		err := h.DeleteThread(microAccountCtx(), &pb.DeleteThreadRequest{
			Id: cRsp.Thread.Id,
		}, &pb.DeleteThreadResponse{})
		assert.NoError(t, err)
	})
}
