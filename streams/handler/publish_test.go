package handler_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/micro/services/streams/handler"
	pb "github.com/micro/services/streams/proto"
	"github.com/stretchr/testify/assert"
)

func TestPublish(t *testing.T) {
	msg := "{\"foo\":\"bar\"}"
	topic := uuid.New().String()

	t.Run("MissingTopic", func(t *testing.T) {
		h := testHandler(t)
		err := h.Publish(context.TODO(), &pb.Message{Message: msg}, &pb.PublishResponse{})
		assert.Equal(t, handler.ErrMissingTopic, err)
		assert.Zero(t, h.Events.(*eventsMock).PublishCount)
	})

	t.Run("MissingMessage", func(t *testing.T) {
		h := testHandler(t)
		err := h.Publish(context.TODO(), &pb.Message{Topic: topic}, &pb.PublishResponse{})
		assert.Equal(t, handler.ErrMissingMessage, err)
		assert.Zero(t, h.Events.(*eventsMock).PublishCount)
	})

	t.Run("ValidMessage", func(t *testing.T) {
		h := testHandler(t)
		err := h.Publish(context.TODO(), &pb.Message{
			Topic: topic, Message: msg,
		}, &pb.PublishResponse{})
		assert.NoError(t, err)
		assert.Equal(t, 1, h.Events.(*eventsMock).PublishCount)
		assert.Equal(t, msg, h.Events.(*eventsMock).PublishMessage)
		assert.Equal(t, topic, h.Events.(*eventsMock).PublishTopic)
	})
}
