package handler_test

import (
	"testing"

	"github.com/micro/services/users/handler"
	pb "github.com/micro/services/users/proto"
	"github.com/stretchr/testify/assert"
)

func TestUpdate(t *testing.T) {
	h := testHandler(t)

	t.Run("MissingID", func(t *testing.T) {
		var rsp pb.UpdateResponse
		err := h.Update(microAccountCtx(), &pb.UpdateRequest{}, &rsp)
		assert.Equal(t, handler.ErrMissingID, err)
		assert.Nil(t, rsp.User)
	})

	t.Run("NotFound", func(t *testing.T) {
		var rsp pb.UpdateResponse
		err := h.Update(microAccountCtx(), &pb.UpdateRequest{Id: "foo"}, &rsp)
		assert.Equal(t, handler.ErrNotFound, err)
		assert.Nil(t, rsp.User)
	})

	// create some mock data
	var cRsp1 pb.CreateResponse
	cReq1 := pb.CreateRequest{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john@doe.com",
		Password:  "passwordabc",
	}
	err := h.Create(microAccountCtx(), &cReq1, &cRsp1)
	assert.NoError(t, err)
	if cRsp1.User == nil {
		t.Fatal("No user returned")
		return
	}

	var cRsp2 pb.CreateResponse
	cReq2 := pb.CreateRequest{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "johndoe@gmail.com",
		Password:  "passwordabc",
	}
	err = h.Create(microAccountCtx(), &cReq2, &cRsp2)
	assert.NoError(t, err)
	if cRsp2.User == nil {
		t.Fatal("No user returned")
		return
	}

	t.Run("BlankFirstName", func(t *testing.T) {
		var rsp pb.UpdateResponse
		err := h.Update(microAccountCtx(), &pb.UpdateRequest{
			Id: cRsp1.User.Id, FirstName: "",
		}, &rsp)
		assert.Equal(t, handler.ErrMissingFirstName, err)
		assert.Nil(t, rsp.User)
	})

	t.Run("BlankLastName", func(t *testing.T) {
		var rsp pb.UpdateResponse
		err := h.Update(microAccountCtx(), &pb.UpdateRequest{
			Id: cRsp1.User.Id, LastName: "",
		}, &rsp)
		assert.Equal(t, handler.ErrMissingLastName, err)
		assert.Nil(t, rsp.User)
	})

	t.Run("BlankLastName", func(t *testing.T) {
		var rsp pb.UpdateResponse
		err := h.Update(microAccountCtx(), &pb.UpdateRequest{
			Id: cRsp1.User.Id, LastName: "",
		}, &rsp)
		assert.Equal(t, handler.ErrMissingLastName, err)
		assert.Nil(t, rsp.User)
	})

	t.Run("BlankEmail", func(t *testing.T) {
		var rsp pb.UpdateResponse
		err := h.Update(microAccountCtx(), &pb.UpdateRequest{
			Id: cRsp1.User.Id, Email: "",
		}, &rsp)
		assert.Equal(t, handler.ErrMissingEmail, err)
		assert.Nil(t, rsp.User)
	})

	t.Run("InvalidEmail", func(t *testing.T) {
		var rsp pb.UpdateResponse
		err := h.Update(microAccountCtx(), &pb.UpdateRequest{
			Id: cRsp1.User.Id, Email: "",
		}, &rsp)
		assert.Equal(t, handler.ErrInvalidEmail, err)
		assert.Nil(t, rsp.User)
	})

	t.Run("EmailAlreadyExists", func(t *testing.T) {
		var rsp pb.UpdateResponse
		err := h.Update(microAccountCtx(), &pb.UpdateRequest{
			Id: cRsp1.User.Id, Email: cRsp2.User.Email,
		}, &rsp)
		assert.Equal(t, handler.ErrDuplicateEmail, err)
		assert.Nil(t, rsp.User)
	})

	t.Run("Valid", func(t *testing.T) {
		uReq := pb.UpdateRequest{
			Id:        cRsp1.User.Id,
			Email:     "foobar@gmail.com",
			FirstName: "Foo",
			LastName:  "Bar",
		}
		var uRsp pb.UpdateResponse
		err := h.Update(microAccountCtx(), &uReq, &uRsp)
		assert.NoError(t, err)
		if uRsp.User == nil {
			t.Error("No user returned")
			return
		}
		assert.Equal(t, cRsp1.User.Id, uRsp.User.Id)
		assert.Equal(t, uReq.Email, uRsp.User.Email)
		assert.Equal(t, uReq.FirstName, uRsp.User.FirstName)
		assert.Equal(t, uReq.LastName, uRsp.User.LastName)
	})

	t.Run("UpdatePassword", func(t *testing.T) {
		uReq := pb.UpdateRequest{
			Id:       cRsp2.User.Id,
			Password: "helloworld",
		}
		err := h.Update(microAccountCtx(), &uReq, &pb.UpdateResponse{})
		assert.NoError(t, err)

		lReq := pb.LoginRequest{
			Email:    cRsp2.User.Email,
			Password: "helloworld",
		}
		err = h.Login(microAccountCtx(), &lReq, &pb.LoginResponse{})
		assert.NoError(t, err)
	})
}
