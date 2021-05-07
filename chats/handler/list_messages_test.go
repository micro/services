package handler_test

import (
	"sort"
	"strconv"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/micro/services/chats/handler"
	pb "github.com/micro/services/chats/proto"
	"github.com/stretchr/testify/assert"
)

func TestListMessages(t *testing.T) {
	h := testHandler(t)
	h.Time = time.Now

	// seed some data
	var chatRsp pb.CreateChatResponse
	err := h.CreateChat(microAccountCtx(), &pb.CreateChatRequest{
		UserIds: []string{uuid.New().String(), uuid.New().String()},
	}, &chatRsp)
	assert.NoError(t, err)
	if err != nil {
		return
	}

	msgs := make([]*pb.Message, 50)
	for i := 0; i < len(msgs); i++ {
		var rsp pb.SendMessageResponse
		err := h.SendMessage(microAccountCtx(), &pb.SendMessageRequest{
			ChatId:   chatRsp.Chat.Id,
			AuthorId: uuid.New().String(),
			Text:     strconv.Itoa(i),
		}, &rsp)
		assert.NoError(t, err)
		msgs[i] = rsp.Message
	}

	t.Run("MissingChatID", func(t *testing.T) {
		var rsp pb.ListMessagesResponse
		err := h.ListMessages(microAccountCtx(), &pb.ListMessagesRequest{}, &rsp)
		assert.Equal(t, handler.ErrMissingChatID, err)
		assert.Nil(t, rsp.Messages)
	})

	t.Run("NoOffset", func(t *testing.T) {
		var rsp pb.ListMessagesResponse
		err := h.ListMessages(microAccountCtx(), &pb.ListMessagesRequest{
			ChatId: chatRsp.Chat.Id,
		}, &rsp)
		assert.NoError(t, err)

		if len(rsp.Messages) != handler.DefaultLimit {
			t.Fatalf("Expected %v messages but got %v", handler.DefaultLimit, len(rsp.Messages))
			return
		}
		expected := msgs[25:]
		sortMessages(rsp.Messages)
		for i, msg := range rsp.Messages {
			assertMessagesMatch(t, expected[i], msg)
		}
	})

	t.Run("LimitSet", func(t *testing.T) {
		var rsp pb.ListMessagesResponse
		err := h.ListMessages(microAccountCtx(), &pb.ListMessagesRequest{
			ChatId: chatRsp.Chat.Id,
			Limit:  10,
		}, &rsp)
		assert.NoError(t, err)

		if len(rsp.Messages) != 10 {
			t.Fatalf("Expected %v messages but got %v", 10, len(rsp.Messages))
			return
		}
		expected := msgs[40:]
		sortMessages(rsp.Messages)
		for i, msg := range rsp.Messages {
			assertMessagesMatch(t, expected[i], msg)
		}
	})

	t.Run("OffsetAndLimit", func(t *testing.T) {
		var rsp pb.ListMessagesResponse
		err := h.ListMessages(microAccountCtx(), &pb.ListMessagesRequest{
			ChatId: chatRsp.Chat.Id,
			Limit:  5,
			Offset: 15,
		}, &rsp)
		assert.NoError(t, err)

		if len(rsp.Messages) != 5 {
			t.Fatalf("Expected %v messages but got %v", 5, len(rsp.Messages))
			return
		}
		expected := msgs[30:35]
		sortMessages(rsp.Messages)
		for i, msg := range rsp.Messages {
			assertMessagesMatch(t, expected[i], msg)
		}
	})
}

// sortMessages by the time they were sent
func sortMessages(msgs []*pb.Message) {
	sort.Slice(msgs, func(i, j int) bool {
		if msgs[i].SentAt == 0 || msgs[j].SentAt == 0 {
			return true
		}
		return msgs[i].SentAt < msgs[j].SentAt
	})
}
