package handler_test

import (
	"testing"

	pb "github.com/micro/services/users/proto"
	"github.com/stretchr/testify/assert"
)

func TestList(t *testing.T) {
	h := testHandler(t)

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

	var rsp pb.ListResponse
	err = h.List(microAccountCtx(), &pb.ListRequest{}, &rsp)
	assert.NoError(t, err)
	if rsp.Users == nil {
		t.Error("No users returned")
		return
	}

	var u1Found, u2Found bool
	for _, u := range rsp.Users {
		switch u.Id {
		case cRsp1.User.Id:
			assertUsersMatch(t, cRsp1.User, u)
			u1Found = true
		case cRsp2.User.Id:
			assertUsersMatch(t, cRsp2.User, u)
			u2Found = true
		default:
			t.Fatal("Unexpected user returned")
			return
		}
	}
	assert.True(t, u1Found)
	assert.True(t, u2Found)
}
