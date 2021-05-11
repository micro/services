package handler_test

import (
	"testing"

	"github.com/micro/services/threads/handler"
	pb "github.com/micro/services/threads/proto"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestCreateThread(t *testing.T) {
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
			var rsp pb.CreateThreadResponse
			err := h.CreateThread(microAccountCtx(), &pb.CreateThreadRequest{
				Topic: tc.Topic, GroupId: tc.GroupID,
			}, &rsp)

			assert.Equal(t, tc.Error, err)
			if tc.Error != nil {
				assert.Nil(t, rsp.Thread)
				return
			}

			assertThreadsMatch(t, &pb.Thread{
				CreatedAt: handler.FormatTime(h.Time()),
				GroupId:   tc.GroupID,
				Topic:     tc.Topic,
			}, rsp.Thread)
		})
	}
}
