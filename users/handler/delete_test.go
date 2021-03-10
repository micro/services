package handler_test

import (
	"context"
	"testing"

	"github.com/micro/services/users/handler"
	pb "github.com/micro/services/users/proto"
	"github.com/stretchr/testify/assert"
)

func TestDelete(t *testing.T) {
	h := testHandler(t)

	t.Run("MissingID", func(t *testing.T) {
		err := h.Delete(context.TODO(), &pb.DeleteRequest{}, &pb.DeleteResponse{})
		assert.Equal(t, handler.ErrMissingID, err)
	})

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

	t.Run("Valid", func(t *testing.T) {
		err := h.Delete(context.TODO(), &pb.DeleteRequest{
			Id: cRsp.User.Id,
		}, &pb.DeleteResponse{})
		assert.NoError(t, err)

		// check it was actually deleted
		var rsp pb.ReadResponse
		err = h.Read(context.TODO(), &pb.ReadRequest{
			Ids: []string{cRsp.User.Id},
		}, &rsp)
		assert.NoError(t, err)
		assert.Nil(t, rsp.Users[cRsp.User.Id])
	})

	t.Run("Retry", func(t *testing.T) {
		err := h.Delete(context.TODO(), &pb.DeleteRequest{
			Id: cRsp.User.Id,
		}, &pb.DeleteResponse{})
		assert.NoError(t, err)
	})
}
