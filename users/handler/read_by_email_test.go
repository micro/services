package handler_test

import (
	"strings"
	"testing"

	"github.com/micro/services/users/handler"
	pb "github.com/micro/services/users/proto"
	"github.com/stretchr/testify/assert"
)

func TestReadByEmail(t *testing.T) {
	h := testHandler(t)

	t.Run("MissingEmails", func(t *testing.T) {
		var rsp pb.ReadByEmailResponse
		err := h.ReadByEmail(microAccountCtx(), &pb.ReadByEmailRequest{}, &rsp)
		assert.Equal(t, handler.ErrMissingEmails, err)
		assert.Nil(t, rsp.Users)
	})

	t.Run("NotFound", func(t *testing.T) {
		var rsp pb.ReadByEmailResponse
		err := h.ReadByEmail(microAccountCtx(), &pb.ReadByEmailRequest{Emails: []string{"foo"}}, &rsp)
		assert.Nil(t, err)
		if rsp.Users == nil {
			t.Fatal("Expected the users object to not be nil")
		}
		assert.Nil(t, rsp.Users["foo"])
	})

	// create some mock data
	var rsp1 pb.CreateResponse
	req1 := pb.CreateRequest{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john@doe.com",
		Password:  "passwordabc",
	}
	err := h.Create(microAccountCtx(), &req1, &rsp1)
	assert.NoError(t, err)
	if rsp1.User == nil {
		t.Fatal("No user returned")
		return
	}

	var rsp2 pb.CreateResponse
	req2 := pb.CreateRequest{
		FirstName: "Apple",
		LastName:  "Tree",
		Email:     "apple@tree.com",
		Password:  "passwordabc",
	}
	err = h.Create(microAccountCtx(), &req2, &rsp2)
	assert.NoError(t, err)
	if rsp2.User == nil {
		t.Fatal("No user returned")
		return
	}

	// test the read
	var rsp pb.ReadByEmailResponse
	err = h.ReadByEmail(microAccountCtx(), &pb.ReadByEmailRequest{
		Emails: []string{rsp1.User.Email, strings.ToUpper(rsp2.User.Email)},
	}, &rsp)
	assert.NoError(t, err)

	if rsp.Users == nil {
		t.Fatal("Users not returned")
		return
	}
	assert.NotNil(t, rsp.Users[rsp1.User.Email])
	assert.NotNil(t, rsp.Users[rsp2.User.Email])

	// check the users match
	if u := rsp.Users[rsp1.User.Email]; u != nil {
		assert.Equal(t, rsp1.User.Id, u.Id)
		assert.Equal(t, rsp1.User.FirstName, u.FirstName)
		assert.Equal(t, rsp1.User.LastName, u.LastName)
		assert.Equal(t, rsp1.User.Email, u.Email)
	}
	if u := rsp.Users[rsp2.User.Email]; u != nil {
		assert.Equal(t, rsp2.User.Id, u.Id)
		assert.Equal(t, rsp2.User.FirstName, u.FirstName)
		assert.Equal(t, rsp2.User.LastName, u.LastName)
		assert.Equal(t, rsp2.User.Email, u.Email)
	}
}
