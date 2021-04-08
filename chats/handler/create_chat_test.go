package handler_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/micro/services/chats/handler"
	pb "github.com/micro/services/chats/proto"
	"github.com/stretchr/testify/assert"
)

func TestCreateChat(t *testing.T) {
	userIDs := []string{uuid.New().String(), uuid.New().String()}

	tt := []struct {
		Name    string
		UserIDs []string
		Error   error
	}{
		{
			Name:  "NoUserIDs",
			Error: handler.ErrMissingUserIDs,
		},
		{
			Name:    "OneUserID",
			UserIDs: userIDs[1:],
			Error:   handler.ErrMissingUserIDs,
		},
		{
			Name:    "Valid",
			UserIDs: userIDs,
		},
		{
			Name:    "Repeat",
			UserIDs: userIDs,
		},
	}

	var chat *pb.Chat
	h := testHandler(t)

	for _, tc := range tt {
		t.Run(tc.Name, func(t *testing.T) {
			var rsp pb.CreateChatResponse
			err := h.CreateChat(microAccountCtx(), &pb.CreateChatRequest{
				UserIds: tc.UserIDs,
			}, &rsp)

			assert.Equal(t, tc.Error, err)
			if tc.Error != nil {
				return
			}

			assert.NotNil(t, rsp.Chat)
			if chat == nil {
				chat = rsp.Chat
			} else {
				assertChatsMatch(t, chat, rsp.Chat)
			}
		})
	}
}
