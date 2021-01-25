package handler_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/micro/services/streams/handler"
	pb "github.com/micro/services/streams/proto"
	"github.com/stretchr/testify/assert"
)

func TestListConversations(t *testing.T) {
	h := testHandler(t)

	// seed some data
	var cRsp1 pb.CreateConversationResponse
	err := h.CreateConversation(context.TODO(), &pb.CreateConversationRequest{
		Topic: "HelloWorld", GroupId: uuid.New().String(),
	}, &cRsp1)
	if err != nil {
		t.Fatalf("Error creating conversation: %v", err)
		return
	}
	var cRsp2 pb.CreateConversationResponse
	err = h.CreateConversation(context.TODO(), &pb.CreateConversationRequest{
		Topic: "FooBar", GroupId: uuid.New().String(),
	}, &cRsp2)
	if err != nil {
		t.Fatalf("Error creating conversation: %v", err)
		return
	}

	t.Run("MissingGroupID", func(t *testing.T) {
		var rsp pb.ListConversationsResponse
		err := h.ListConversations(context.TODO(), &pb.ListConversationsRequest{}, &rsp)
		assert.Equal(t, handler.ErrMissingGroupID, err)
		assert.Nil(t, rsp.Conversations)
	})

	t.Run("Valid", func(t *testing.T) {
		var rsp pb.ListConversationsResponse
		err := h.ListConversations(context.TODO(), &pb.ListConversationsRequest{
			GroupId: cRsp1.Conversation.GroupId,
		}, &rsp)

		assert.NoError(t, err)
		if len(rsp.Conversations) != 1 {
			t.Fatalf("Expected 1 conversation to be returned, got %v", len(rsp.Conversations))
			return
		}

		assertConversationsMatch(t, cRsp1.Conversation, rsp.Conversations[0])
	})
}
