package handler_test

import (
	"context"
	"sort"
	"strconv"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/micro/services/chats/handler"
	pb "github.com/micro/services/chats/proto"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func TestListMessages(t *testing.T) {
	h := testHandler(t)
	h.Time = time.Now

	// seed some data
	var chatRsp pb.CreateChatResponse
	err := h.CreateChat(context.TODO(), &pb.CreateChatRequest{
		UserIds: []string{uuid.New().String(), uuid.New().String()},
	}, &chatRsp)
	assert.NoError(t, err)
	if err != nil {
		return
	}

	msgs := make([]*pb.Message, 50)
	for i := 0; i < len(msgs); i++ {
		var rsp pb.CreateMessageResponse
		err := h.CreateMessage(context.TODO(), &pb.CreateMessageRequest{
			ChatId:   chatRsp.Chat.Id,
			AuthorId: uuid.New().String(),
			Text:     strconv.Itoa(i),
		}, &rsp)
		assert.NoError(t, err)
		msgs[i] = rsp.Message
	}

	t.Run("MissingChatID", func(t *testing.T) {
		var rsp pb.ListMessagesResponse
		err := h.ListMessages(context.TODO(), &pb.ListMessagesRequest{}, &rsp)
		assert.Equal(t, handler.ErrMissingChatID, err)
		assert.Nil(t, rsp.Messages)
	})

	t.Run("NoOffset", func(t *testing.T) {
		var rsp pb.ListMessagesResponse
		err := h.ListMessages(context.TODO(), &pb.ListMessagesRequest{
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
		err := h.ListMessages(context.TODO(), &pb.ListMessagesRequest{
			ChatId: chatRsp.Chat.Id,
			Limit:  &wrapperspb.Int32Value{Value: 10},
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
		err := h.ListMessages(context.TODO(), &pb.ListMessagesRequest{
			ChatId:     chatRsp.Chat.Id,
			Limit:      &wrapperspb.Int32Value{Value: 5},
			SentBefore: msgs[20].SentAt,
		}, &rsp)
		assert.NoError(t, err)

		if len(rsp.Messages) != 5 {
			t.Fatalf("Expected %v messages but got %v", 5, len(rsp.Messages))
			return
		}
		expected := msgs[15:20]
		sortMessages(rsp.Messages)
		for i, msg := range rsp.Messages {
			assertMessagesMatch(t, expected[i], msg)
		}
	})
}

// sortMessages by the time they were sent
func sortMessages(msgs []*pb.Message) {
	sort.Slice(msgs, func(i, j int) bool {
		if msgs[i].SentAt == nil || msgs[j].SentAt == nil {
			return true
		}
		return msgs[i].SentAt.AsTime().Before(msgs[j].SentAt.AsTime())
	})
}
