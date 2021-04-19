package handler_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/micro/services/threads/handler"
	pb "github.com/micro/services/threads/proto"
	"github.com/stretchr/testify/assert"
)

func TestRecentMessages(t *testing.T) {
	h := testHandler(t)
	h.Time = time.Now

	// seed some data
	ids := make([]string, 3)
	convos := make(map[string][]*pb.Message, 3)
	for i := 0; i < 3; i++ {
		var convRsp pb.CreateConversationResponse
		err := h.CreateConversation(microAccountCtx(), &pb.CreateConversationRequest{
			Topic: "TestRecentMessages", GroupId: uuid.New().String(),
		}, &convRsp)
		assert.NoError(t, err)
		if err != nil {
			return
		}

		convos[convRsp.Conversation.Id] = make([]*pb.Message, 50)
		ids[i] = convRsp.Conversation.Id

		for j := 0; j < 50; j++ {
			var rsp pb.CreateMessageResponse
			err := h.CreateMessage(microAccountCtx(), &pb.CreateMessageRequest{
				ConversationId: convRsp.Conversation.Id,
				AuthorId:       uuid.New().String(),
				Text:           fmt.Sprintf("Conversation %v, Message %v", i, j),
			}, &rsp)
			assert.NoError(t, err)
			convos[convRsp.Conversation.Id][j] = rsp.Message
		}
	}

	t.Run("MissingConversationIDs", func(t *testing.T) {
		var rsp pb.RecentMessagesResponse
		err := h.RecentMessages(microAccountCtx(), &pb.RecentMessagesRequest{}, &rsp)
		assert.Equal(t, handler.ErrMissingConversationIDs, err)
		assert.Nil(t, rsp.Messages)
	})

	t.Run("LimitSet", func(t *testing.T) {
		var rsp pb.RecentMessagesResponse
		err := h.RecentMessages(microAccountCtx(), &pb.RecentMessagesRequest{
			ConversationIds:      ids,
			LimitPerConversation: 10,
		}, &rsp)
		assert.NoError(t, err)

		if len(rsp.Messages) != 30 {
			t.Fatalf("Expected %v messages but got %v", 30, len(rsp.Messages))
			return
		}
		var expected []*pb.Message
		for _, msgs := range convos {
			expected = append(expected, msgs[40:]...)
		}
		sortMessages(expected)
		sortMessages(rsp.Messages)
		for _, msg := range rsp.Messages {
			assertMessagesMatch(t, getMsg(msg.Id, expected), msg)
		}
	})

	t.Run("NoLimitSet", func(t *testing.T) {
		reducedIDs := ids[:2]

		var rsp pb.RecentMessagesResponse
		err := h.RecentMessages(microAccountCtx(), &pb.RecentMessagesRequest{
			ConversationIds: reducedIDs,
		}, &rsp)
		assert.NoError(t, err)

		if len(rsp.Messages) != 50 {
			t.Fatalf("Expected %v messages but got %v", 50, len(rsp.Messages))
			return
		}
		var expected []*pb.Message
		for _, id := range reducedIDs {
			expected = append(expected, convos[id][25:]...)
		}
		sortMessages(expected)
		sortMessages(rsp.Messages)
		for _, msg := range rsp.Messages {
			assertMessagesMatch(t, getMsg(msg.Id, expected), msg)
		}
	})
}
