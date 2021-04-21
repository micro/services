package handler_test

import (
	"context"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/micro/micro/v3/service/auth"
	"github.com/micro/services/streams/handler"
	pb "github.com/micro/services/streams/proto"
	"github.com/stretchr/testify/assert"
)

func TestPublish(t *testing.T) {
	msg := "{\"foo\":\"bar\"}"
	topic := strings.ReplaceAll(uuid.New().String(), "-", "")

	t.Run("MissingTopic", func(t *testing.T) {
		h := testHandler(t)
		ctx := auth.ContextWithAccount(context.TODO(), &auth.Account{Issuer: "foo"})
		err := h.Publish(ctx, &pb.Message{Message: msg}, &pb.PublishResponse{})
		assert.Equal(t, handler.ErrMissingTopic, err)
		assert.Zero(t, h.Events.(*eventsMock).PublishCount)
	})

	t.Run("MissingMessage", func(t *testing.T) {
		h := testHandler(t)
		ctx := auth.ContextWithAccount(context.TODO(), &auth.Account{Issuer: "foo"})
		err := h.Publish(ctx, &pb.Message{Topic: topic}, &pb.PublishResponse{})
		assert.Equal(t, handler.ErrMissingMessage, err)
		assert.Zero(t, h.Events.(*eventsMock).PublishCount)
	})

	t.Run("ValidMessage", func(t *testing.T) {
		h := testHandler(t)
		ctx := auth.ContextWithAccount(context.TODO(), &auth.Account{Issuer: "foo", ID: "foo-id"})
		err := h.Publish(ctx, &pb.Message{
			Topic: topic, Message: msg,
		}, &pb.PublishResponse{})
		assert.NoError(t, err)
		assert.Equal(t, 1, h.Events.(*eventsMock).PublishCount)
		assert.Equal(t, msg, h.Events.(*eventsMock).PublishMessage)
		// topic is prefixed with acc issuer to implement multitenancy
		assert.Equal(t, "foo.foo-id."+topic, h.Events.(*eventsMock).PublishTopic)
	})
}
