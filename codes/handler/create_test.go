package handler_test

import (
	"context"
	"testing"
	"time"

	"github.com/micro/micro/v3/service/auth"
	"github.com/micro/services/codes/handler"
	pb "github.com/micro/services/codes/proto"
	"github.com/stretchr/testify/assert"
)

func TestCreate(t *testing.T) {
	h := testHandler(t)

	t.Run("MissingIdentity", func(t *testing.T) {
		var rsp pb.CreateResponse
		err := h.Create(microAccountCtx(), &pb.CreateRequest{}, &rsp)
		assert.Equal(t, handler.ErrMissingIdentity, err)
		assert.Empty(t, rsp.Code)
	})

	t.Run("NoExpiry", func(t *testing.T) {
		var rsp pb.CreateResponse
		err := h.Create(microAccountCtx(), &pb.CreateRequest{Identity: "07503196715"}, &rsp)
		assert.NoError(t, err)
		assert.NotEmpty(t, rsp.Code)
	})

	t.Run("WithExpiry", func(t *testing.T) {
		var rsp pb.CreateResponse
		err := h.Create(microAccountCtx(), &pb.CreateRequest{
			Identity:  "demo@m3o.com",
			ExpiresAt: time.Now().Unix(),
		}, &rsp)
		assert.NoError(t, err)
		assert.NotEmpty(t, rsp.Code)
	})
}

func microAccountCtx() context.Context {
	return auth.ContextWithAccount(context.TODO(), &auth.Account{
		Issuer: "micro",
	})
}
