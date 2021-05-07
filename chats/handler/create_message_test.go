package handler_test

import (
	"testing"
	"time"

	"github.com/micro/services/chats/handler"
	pb "github.com/micro/services/chats/proto"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestSendMessage(t *testing.T) {
	h := testHandler(t)

	// seed some data
	var cRsp pb.CreateChatResponse
	err := h.CreateChat(microAccountCtx(), &pb.CreateChatRequest{
		UserIds: []string{uuid.New().String(), uuid.New().String()},
	}, &cRsp)
	if err != nil {
		t.Fatalf("Error creating chat: %v", err)
		return
	}

	iid := uuid.New().String()
	tt := []struct {
		Name     string
		AuthorID string
		ChatID   string
		Text     string
		Error    error
		ID       string
	}{
		{
			Name:     "MissingChatID",
			Text:     "HelloWorld",
			AuthorID: uuid.New().String(),
			Error:    handler.ErrMissingChatID,
		},
		{
			Name:   "MissingAuthorID",
			ChatID: uuid.New().String(),
			Text:   "HelloWorld",
			Error:  handler.ErrMissingAuthorID,
		},
		{
			Name:     "MissingText",
			ChatID:   uuid.New().String(),
			AuthorID: uuid.New().String(),
			Error:    handler.ErrMissingText,
		},
		{
			Name:     "ChatNotFound",
			ChatID:   uuid.New().String(),
			AuthorID: uuid.New().String(),
			Text:     "HelloWorld",
			Error:    handler.ErrNotFound,
		},
		{
			Name:     "WithoutID",
			ChatID:   cRsp.Chat.Id,
			AuthorID: uuid.New().String(),
			Text:     "HelloWorld",
		},
		{
			Name:     "WithID",
			ChatID:   cRsp.Chat.Id,
			AuthorID: "johndoe",
			Text:     "HelloWorld",
			ID:       iid,
		},
		{
			Name:     "RepeatID",
			ChatID:   cRsp.Chat.Id,
			AuthorID: "johndoe",
			Text:     "HelloWorld",
			ID:       iid,
		},
	}

	for _, tc := range tt {
		t.Run(tc.Name, func(t *testing.T) {
			var rsp pb.SendMessageResponse
			err := h.SendMessage(microAccountCtx(), &pb.SendMessageRequest{
				AuthorId: tc.AuthorID,
				ChatId:   tc.ChatID,
				Text:     tc.Text,
				Id:       tc.ID,
			}, &rsp)

			assert.Equal(t, tc.Error, err)
			if tc.Error != nil {
				assert.Nil(t, rsp.Message)
				return
			}

			assertMessagesMatch(t, &pb.Message{
				AuthorId: tc.AuthorID,
				ChatId:   tc.ChatID,
				SentAt:   h.Time().Format(time.RFC3339Nano),
				Text:     tc.Text,
				Id:       tc.ID,
			}, rsp.Message)
		})
	}
}
