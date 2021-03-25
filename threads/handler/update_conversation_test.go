package handler_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/micro/services/threads/handler"
	pb "github.com/micro/services/threads/proto"
	"github.com/stretchr/testify/assert"
)

func TestUpdateConversation(t *testing.T) {
	h := testHandler(t)

	// seed some data
	var cRsp pb.CreateConversationResponse
	err := h.CreateConversation(microAccountCtx(), &pb.CreateConversationRequest{
		Topic: "HelloWorld", GroupId: uuid.New().String(),
	}, &cRsp)
	if err != nil {
		t.Fatalf("Error creating conversation: %v", err)
		return
	}

	t.Run("MissingID", func(t *testing.T) {
		err := h.UpdateConversation(microAccountCtx(), &pb.UpdateConversationRequest{
			Topic: "NewTopic",
		}, &pb.UpdateConversationResponse{})
		assert.Equal(t, handler.ErrMissingID, err)
	})

	t.Run("MissingTopic", func(t *testing.T) {
		err := h.UpdateConversation(microAccountCtx(), &pb.UpdateConversationRequest{
			Id: uuid.New().String(),
		}, &pb.UpdateConversationResponse{})
		assert.Equal(t, handler.ErrMissingTopic, err)
	})

	t.Run("InvalidID", func(t *testing.T) {
		err := h.UpdateConversation(microAccountCtx(), &pb.UpdateConversationRequest{
			Id:    uuid.New().String(),
			Topic: "NewTopic",
		}, &pb.UpdateConversationResponse{})
		assert.Equal(t, handler.ErrNotFound, err)
	})

	t.Run("Valid", func(t *testing.T) {
		err := h.UpdateConversation(microAccountCtx(), &pb.UpdateConversationRequest{
			Id:    cRsp.Conversation.Id,
			Topic: "NewTopic",
		}, &pb.UpdateConversationResponse{})
		assert.NoError(t, err)

		var rsp pb.ReadConversationResponse
		err = h.ReadConversation(microAccountCtx(), &pb.ReadConversationRequest{
			Id: cRsp.Conversation.Id,
		}, &rsp)
		assert.NoError(t, err)
		if rsp.Conversation == nil {
			t.Fatal("No conversation returned")
			return
		}
		assert.Equal(t, "NewTopic", rsp.Conversation.Topic)
	})
}
