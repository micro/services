package handler_test

import (
	"context"
	"testing"

	pb "github.com/micro/services/streams/proto"
	"github.com/stretchr/testify/assert"
)

func TestToken(t *testing.T) {
	h := testHandler(t)

	t.Run("WithoutTopic", func(t *testing.T) {
		var rsp pb.TokenResponse
		err := h.Token(context.TODO(), &pb.TokenRequest{}, &rsp)
		assert.NoError(t, err)
		assert.NotEmpty(t, rsp.Token)
	})

	t.Run("WithTopic", func(t *testing.T) {
		var rsp pb.TokenResponse
		err := h.Token(context.TODO(), &pb.TokenRequest{Topic: "helloworld"}, &rsp)
		assert.NoError(t, err)
		assert.NotEmpty(t, rsp.Token)
	})
}
