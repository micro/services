package handler_test

import (
	"testing"

	"github.com/micro/services/threads/handler"
	pb "github.com/micro/services/threads/proto"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestCreateMessage(t *testing.T) {
	h := testHandler(t)

	// seed some data
	var cRsp pb.CreateThreadResponse
	err := h.CreateThread(microAccountCtx(), &pb.CreateThreadRequest{
		Topic: "HelloWorld", GroupId: uuid.New().String(),
	}, &cRsp)
	if err != nil {
		t.Fatalf("Error creating thread: %v", err)
		return
	}

	iid := uuid.New().String()
	tt := []struct {
		Name     string
		AuthorID string
		ThreadID string
		ID       string
		Text     string
		Error    error
	}{
		{
			Name:     "MissingThreadID",
			Text:     "HelloWorld",
			AuthorID: uuid.New().String(),
			Error:    handler.ErrMissingThreadID,
		},
		{
			Name:     "MissingAuthorID",
			ThreadID: uuid.New().String(),
			Text:     "HelloWorld",
			Error:    handler.ErrMissingAuthorID,
		},
		{
			Name:     "MissingText",
			ThreadID: uuid.New().String(),
			AuthorID: uuid.New().String(),
			Error:    handler.ErrMissingText,
		},
		{
			Name:     "ThreadNotFound",
			ThreadID: uuid.New().String(),
			AuthorID: uuid.New().String(),
			Text:     "HelloWorld",
			Error:    handler.ErrNotFound,
		},
		{
			Name:     "NoID",
			ThreadID: cRsp.Thread.Id,
			AuthorID: uuid.New().String(),
			Text:     "HelloWorld",
		},
		{
			Name:     "WithID",
			ThreadID: cRsp.Thread.Id,
			Text:     "HelloWorld",
			AuthorID: "johndoe",
			ID:       iid,
		},
		{
			Name:     "RepeatID",
			ThreadID: cRsp.Thread.Id,
			Text:     "HelloWorld",
			AuthorID: "johndoe",
			ID:       iid,
		},
	}

	for _, tc := range tt {
		t.Run(tc.Name, func(t *testing.T) {
			var rsp pb.CreateMessageResponse
			err := h.CreateMessage(microAccountCtx(), &pb.CreateMessageRequest{
				AuthorId: tc.AuthorID,
				ThreadId: tc.ThreadID,
				Text:     tc.Text,
				Id:       tc.ID,
			}, &rsp)

			assert.Equal(t, tc.Error, err)
			if tc.Error != nil {
				assert.Nil(t, rsp.Message)
				return
			}

			assertMessagesMatch(t, &pb.Message{
				Id:       tc.ID,
				AuthorId: tc.AuthorID,
				ThreadId: tc.ThreadID,
				SentAt:   handler.FormatTime(h.Time()),
				Text:     tc.Text,
			}, rsp.Message)
		})
	}
}
