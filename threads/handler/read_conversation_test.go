package handler_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/micro/services/threads/handler"
	pb "github.com/micro/services/threads/proto"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func TestReadConversation(t *testing.T) {
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

	tt := []struct {
		Name    string
		ID      string
		GroupID *wrapperspb.StringValue
		Error   error
		Result  *pb.Conversation
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
			ID:     cRsp.Conversation.Id,
			Result: cRsp.Conversation,
		},
		{
			Name:    "IncorrectGroupID",
			ID:      cRsp.Conversation.Id,
			Error:   handler.ErrNotFound,
			GroupID: &wrapperspb.StringValue{Value: uuid.New().String()},
		},
	}
	for _, tc := range tt {
		t.Run(tc.Name, func(t *testing.T) {
			var rsp pb.ReadConversationResponse
			err := h.ReadConversation(context.TODO(), &pb.ReadConversationRequest{
				Id: tc.ID, GroupId: tc.GroupID,
			}, &rsp)
			assert.Equal(t, tc.Error, err)

			if tc.Result == nil {
				assert.Nil(t, rsp.Conversation)
			} else {
				assertConversationsMatch(t, tc.Result, rsp.Conversation)
			}
		})
	}
}
