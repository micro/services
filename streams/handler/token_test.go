package handler_test

import (
	"context"
	"testing"

	"github.com/micro/micro/v3/service/auth"
	pb "github.com/micro/services/streams/proto"
	"github.com/stretchr/testify/assert"
)

func TestToken(t *testing.T) {
	h := testHandler(t)

	t.Run("WithoutTopic", func(t *testing.T) {
		var rsp pb.TokenResponse
		ctx := auth.ContextWithAccount(context.TODO(), &auth.Account{Issuer: "foo"})
		err := h.Token(ctx, &pb.TokenRequest{}, &rsp)
		assert.NoError(t, err)
		assert.NotEmpty(t, rsp.Token)
	})

	t.Run("WithTopic", func(t *testing.T) {
		var rsp pb.TokenResponse
		ctx := auth.ContextWithAccount(context.TODO(), &auth.Account{Issuer: "foo"})
		err := h.Token(ctx, &pb.TokenRequest{Topic: "helloworld"}, &rsp)
		assert.NoError(t, err)
		assert.NotEmpty(t, rsp.Token)
	})
}
