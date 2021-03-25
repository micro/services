package handler_test

import (
	"testing"

	"github.com/micro/services/threads/handler"
	pb "github.com/micro/services/threads/proto"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestDeleteConversation(t *testing.T) {
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
		err := h.DeleteConversation(microAccountCtx(), &pb.DeleteConversationRequest{}, &pb.DeleteConversationResponse{})
		assert.Equal(t, handler.ErrMissingID, err)
	})

	t.Run("Valid", func(t *testing.T) {
		err := h.DeleteConversation(microAccountCtx(), &pb.DeleteConversationRequest{
			Id: cRsp.Conversation.Id,
		}, &pb.DeleteConversationResponse{})
		assert.NoError(t, err)

		err = h.ReadConversation(microAccountCtx(), &pb.ReadConversationRequest{
			Id: cRsp.Conversation.Id,
		}, &pb.ReadConversationResponse{})
		assert.Equal(t, handler.ErrNotFound, err)
	})

	t.Run("Retry", func(t *testing.T) {
		err := h.DeleteConversation(microAccountCtx(), &pb.DeleteConversationRequest{
			Id: cRsp.Conversation.Id,
		}, &pb.DeleteConversationResponse{})
		assert.NoError(t, err)
	})
}
