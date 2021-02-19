package handler_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/micro/services/users/handler"
	pb "github.com/micro/services/users/proto"
	"github.com/stretchr/testify/assert"
)

func TestLogout(t *testing.T) {
	h := testHandler(t)

	t.Run("MissingUserID", func(t *testing.T) {
		err := h.Logout(context.TODO(), &pb.LogoutRequest{}, &pb.LogoutResponse{})
		assert.Equal(t, handler.ErrMissingID, err)
	})

	t.Run("UserNotFound", func(t *testing.T) {
		err := h.Logout(context.TODO(), &pb.LogoutRequest{Id: uuid.New().String()}, &pb.LogoutResponse{})
		assert.Equal(t, handler.ErrNotFound, err)
	})

	t.Run("Valid", func(t *testing.T) {
		// create some mock data
		var cRsp pb.CreateResponse
		cReq := pb.CreateRequest{
			FirstName: "John",
			LastName:  "Doe",
			Email:     "john@doe.com",
			Password:  "passwordabc",
		}
		err := h.Create(context.TODO(), &cReq, &cRsp)
		assert.NoError(t, err)
		if cRsp.User == nil {
			t.Fatal("No user returned")
			return
		}

		err = h.Logout(context.TODO(), &pb.LogoutRequest{Id: cRsp.User.Id}, &pb.LogoutResponse{})
		assert.NoError(t, err)

		err = h.Validate(context.TODO(), &pb.ValidateRequest{Token: cRsp.Token}, &pb.ValidateResponse{})
		assert.Error(t, err)
	})
}
