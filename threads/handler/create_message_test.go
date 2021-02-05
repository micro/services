package handler_test

import (
	"context"
	"testing"

	"github.com/micro/services/threads/handler"
	pb "github.com/micro/services/threads/proto"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestCreateMessage(t *testing.T) {
	h := testHandler(t)

	// seed some data
	var cRsp pb.CreateConversationResponse
	err := h.CreateConversation(context.TODO(), &pb.CreateConversationRequest{
		Topic: "HelloWorld", GroupId: uuid.New().String(),
	}, &cRsp)
	if err != nil {
		t.Fatalf("Error creating conversation: %v", err)
		return
	}

	iid := uuid.New().String()
	tt := []struct {
		Name           string
		AuthorID       string
		ConversationID string
		IdempotentID   string
		Text           string
		Error          error
	}{
		{
			Name:     "MissingConversationID",
			Text:     "HelloWorld",
			AuthorID: uuid.New().String(),
			Error:    handler.ErrMissingConversationID,
		},
		{
			Name:           "MissingAuthorID",
			ConversationID: uuid.New().String(),
			Text:           "HelloWorld",
			Error:          handler.ErrMissingAuthorID,
		},
		{
			Name:           "MissingText",
			ConversationID: uuid.New().String(),
			AuthorID:       uuid.New().String(),
			Error:          handler.ErrMissingText,
		},
		{
			Name:           "ConversationNotFound",
			ConversationID: uuid.New().String(),
			AuthorID:       uuid.New().String(),
			Text:           "HelloWorld",
			Error:          handler.ErrNotFound,
		},
		{
			Name:           "NoIdempotentID",
			ConversationID: cRsp.Conversation.Id,
			AuthorID:       uuid.New().String(),
			Text:           "HelloWorld",
		},
		{
			Name:           "WithIdempotentID",
			ConversationID: cRsp.Conversation.Id,
			Text:           "HelloWorld",
			AuthorID:       "johndoe",
			IdempotentID:   iid,
		},
		{
			Name:           "RepeatIdempotentID",
			ConversationID: cRsp.Conversation.Id,
			Text:           "HelloWorld",
			AuthorID:       "johndoe",
			IdempotentID:   iid,
		},
	}

	for _, tc := range tt {
		t.Run(tc.Name, func(t *testing.T) {
			var rsp pb.CreateMessageResponse
			err := h.CreateMessage(context.TODO(), &pb.CreateMessageRequest{
				AuthorId:       tc.AuthorID,
				ConversationId: tc.ConversationID,
				Text:           tc.Text,
				IdempotentId:   tc.IdempotentID,
			}, &rsp)

			assert.Equal(t, tc.Error, err)
			if tc.Error != nil {
				assert.Nil(t, rsp.Message)
				return
			}

			assertMessagesMatch(t, &pb.Message{
				IdempotentId:   tc.IdempotentID,
				AuthorId:       tc.AuthorID,
				ConversationId: tc.ConversationID,
				SentAt:         timestamppb.New(h.Time()),
				Text:           tc.Text,
			}, rsp.Message)
		})
	}
}
