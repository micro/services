package handler_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/micro/services/threads/handler"
	pb "github.com/micro/services/threads/proto"
	"github.com/stretchr/testify/assert"
)

func TestUpdateThread(t *testing.T) {
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
		err := h.UpdateThread(microAccountCtx(), &pb.UpdateThreadRequest{
			Topic: "NewTopic",
		}, &pb.UpdateThreadResponse{})
		assert.Equal(t, handler.ErrMissingID, err)
	})

	t.Run("MissingTopic", func(t *testing.T) {
		err := h.UpdateThread(microAccountCtx(), &pb.UpdateThreadRequest{
			Id: uuid.New().String(),
		}, &pb.UpdateThreadResponse{})
		assert.Equal(t, handler.ErrMissingTopic, err)
	})

	t.Run("InvalidID", func(t *testing.T) {
		err := h.UpdateThread(microAccountCtx(), &pb.UpdateThreadRequest{
			Id:    uuid.New().String(),
			Topic: "NewTopic",
		}, &pb.UpdateThreadResponse{})
		assert.Equal(t, handler.ErrNotFound, err)
	})

	t.Run("Valid", func(t *testing.T) {
		err := h.UpdateThread(microAccountCtx(), &pb.UpdateThreadRequest{
			Id:    cRsp.Thread.Id,
			Topic: "NewTopic",
		}, &pb.UpdateThreadResponse{})
		assert.NoError(t, err)

		var rsp pb.ReadThreadResponse
		err = h.ReadThread(microAccountCtx(), &pb.ReadThreadRequest{
			Id: cRsp.Thread.Id,
		}, &rsp)
		assert.NoError(t, err)
		if rsp.Thread == nil {
			t.Fatal("No thread returned")
			return
		}
		assert.Equal(t, "NewTopic", rsp.Thread.Topic)
	})
}
