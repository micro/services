package handler_test

import (
	"sort"
	"strconv"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/micro/services/threads/handler"
	pb "github.com/micro/services/threads/proto"
	"github.com/stretchr/testify/assert"
)

func TestListMessages(t *testing.T) {
	h := testHandler(t)
	h.Time = time.Now

	// seed some data
	var convRsp pb.CreateConversationResponse
	err := h.CreateConversation(microAccountCtx(), &pb.CreateConversationRequest{
		Topic: "TestListMessages", GroupId: uuid.New().String(),
	}, &convRsp)
	assert.NoError(t, err)
	if err != nil {
		return
	}

	msgs := make([]*pb.Message, 50)
	for i := 0; i < len(msgs); i++ {
		var rsp pb.CreateMessageResponse
		err := h.CreateMessage(microAccountCtx(), &pb.CreateMessageRequest{
			ConversationId: convRsp.Conversation.Id,
			AuthorId:       uuid.New().String(),
			Text:           strconv.Itoa(i),
		}, &rsp)
		assert.NoError(t, err)
		msgs[i] = rsp.Message
	}

	t.Run("MissingConversationID", func(t *testing.T) {
		var rsp pb.ListMessagesResponse
		err := h.ListMessages(microAccountCtx(), &pb.ListMessagesRequest{}, &rsp)
		assert.Equal(t, handler.ErrMissingConversationID, err)
		assert.Nil(t, rsp.Messages)
	})

	t.Run("NoOffset", func(t *testing.T) {
		var rsp pb.ListMessagesResponse
		err := h.ListMessages(microAccountCtx(), &pb.ListMessagesRequest{
			ConversationId: convRsp.Conversation.Id,
		}, &rsp)
		assert.NoError(t, err)

		if len(rsp.Messages) != handler.DefaultLimit {
			t.Fatalf("Expected %v messages but got %v", handler.DefaultLimit, len(rsp.Messages))
			return
		}
		expected := msgs[25:]
		sortMessages(rsp.Messages)
		for _, msg := range rsp.Messages {
			assertMessagesMatch(t, getMsg(msg.Id, expected), msg)
		}
	})

	t.Run("LimitSet", func(t *testing.T) {
		var rsp pb.ListMessagesResponse
		err := h.ListMessages(microAccountCtx(), &pb.ListMessagesRequest{
			ConversationId: convRsp.Conversation.Id,
			Limit:          10,
		}, &rsp)
		assert.NoError(t, err)

		if len(rsp.Messages) != 10 {
			t.Fatalf("Expected %v messages but got %v", 10, len(rsp.Messages))
			return
		}
		expected := msgs[40:]
		sortMessages(rsp.Messages)
		for _, msg := range rsp.Messages {
			assertMessagesMatch(t, getMsg(msg.Id, expected), msg)
		}
	})

	t.Run("OffsetAndLimit", func(t *testing.T) {
		var rsp pb.ListMessagesResponse
		err := h.ListMessages(microAccountCtx(), &pb.ListMessagesRequest{
			ConversationId: convRsp.Conversation.Id,
			Limit:          5,
			SentBefore:     msgs[20].SentAt,
		}, &rsp)
		assert.NoError(t, err)

		if len(rsp.Messages) != 5 {
			t.Fatalf("Expected %v messages but got %v", 5, len(rsp.Messages))
			return
		}
		expected := msgs[15:20]
		sortMessages(rsp.Messages)
		for _, msg := range rsp.Messages {
			assertMessagesMatch(t, getMsg(msg.Id, expected), msg)
		}
	})
}

func getMsg(id string, msgs []*pb.Message) *pb.Message {
	for _, msg := range msgs {
		if id == msg.Id {
			return msg
		}
	}
	return nil
}

// sortMessages by the time they were sent
func sortMessages(msgs []*pb.Message) {
	sort.Slice(msgs, func(i, j int) bool {
		return msgs[i].Id < msgs[j].Id
	})
}
