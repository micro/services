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

func TestCreateConversation(t *testing.T) {
	tt := []struct {
		Name    string
		GroupID string
		Topic   string
		Error   error
	}{
		{
			Name:  "MissingGroupID",
			Topic: "HelloWorld",
			Error: handler.ErrMissingGroupID,
		},
		{
			Name:    "MissingTopic",
			GroupID: uuid.New().String(),
			Error:   handler.ErrMissingTopic,
		},
		{
			Name:    "Valid",
			GroupID: uuid.New().String(),
			Topic:   "HelloWorld",
		},
	}

	h := testHandler(t)
	for _, tc := range tt {
		t.Run(tc.Name, func(t *testing.T) {
			var rsp pb.CreateConversationResponse
			err := h.CreateConversation(context.TODO(), &pb.CreateConversationRequest{
				Topic: tc.Topic, GroupId: tc.GroupID,
			}, &rsp)

			assert.Equal(t, tc.Error, err)
			if tc.Error != nil {
				assert.Nil(t, rsp.Conversation)
				return
			}

			assertConversationsMatch(t, &pb.Conversation{
				CreatedAt: timestamppb.New(h.Time()),
				GroupId:   tc.GroupID,
				Topic:     tc.Topic,
			}, rsp.Conversation)
		})
	}
}
