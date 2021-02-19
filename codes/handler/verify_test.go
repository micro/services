package handler_test

import (
	"context"
	"testing"
	"time"

	"github.com/micro/services/codes/handler"
	pb "github.com/micro/services/codes/proto"
	"github.com/stretchr/testify/assert"
)

func TestVerify(t *testing.T) {
	h := testHandler(t)

	t.Run("MissingIdentity", func(t *testing.T) {
		var rsp pb.VerifyResponse
		err := h.Verify(context.TODO(), &pb.VerifyRequest{Code: "123456"}, &rsp)
		assert.Equal(t, handler.ErrMissingIdentity, err)
	})

	t.Run("MissingCode", func(t *testing.T) {
		var rsp pb.VerifyResponse
		err := h.Verify(context.TODO(), &pb.VerifyRequest{Identity: "demo@m3o.com"}, &rsp)
		assert.Equal(t, handler.ErrMissingCode, err)
	})

	// generate a code to test
	var cRsp pb.CreateResponse
	err := h.Create(context.TODO(), &pb.CreateRequest{Identity: "demo@m3o.com"}, &cRsp)
	assert.NoError(t, err)

	t.Run("IncorrectCode", func(t *testing.T) {
		var rsp pb.VerifyResponse
		err := h.Verify(context.TODO(), &pb.VerifyRequest{Identity: "demo@m3o.com", Code: "12345"}, &rsp)
		assert.Equal(t, handler.ErrInvalidCode, err)
	})

	t.Run("IncorrectEmail", func(t *testing.T) {
		var rsp pb.VerifyResponse
		err := h.Verify(context.TODO(), &pb.VerifyRequest{Identity: "john@m3o.com", Code: cRsp.Code}, &rsp)
		assert.Equal(t, handler.ErrInvalidCode, err)
	})

	t.Run("ExpiredCode", func(t *testing.T) {
		ot := h.Time
		h.Time = func() time.Time { return time.Now().Add(handler.DefaultTTL * 2) }
		defer func() { h.Time = ot }()

		var rsp pb.VerifyResponse
		err := h.Verify(context.TODO(), &pb.VerifyRequest{Identity: "demo@m3o.com", Code: cRsp.Code}, &rsp)
		assert.Equal(t, handler.ErrExpiredCode, err)
	})

	t.Run("ValidCode", func(t *testing.T) {
		var rsp pb.VerifyResponse
		err := h.Verify(context.TODO(), &pb.VerifyRequest{Identity: "demo@m3o.com", Code: cRsp.Code}, &rsp)
		assert.NoError(t, err)
	})
}
