package handler_test

import (
	"testing"

	"github.com/micro/services/users/handler"
	pb "github.com/micro/services/users/proto"
	"github.com/stretchr/testify/assert"
)

func TestLogin(t *testing.T) {
	h := testHandler(t)

	// create some mock data
	var cRsp pb.CreateResponse
	cReq := pb.CreateRequest{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john@doe.com",
		Password:  "passwordabc",
	}
	err := h.Create(microAccountCtx(), &cReq, &cRsp)
	assert.NoError(t, err)
	if cRsp.User == nil {
		t.Fatal("No user returned")
		return
	}

	tt := []struct {
		Name     string
		Email    string
		Password string
		Error    error
		User     *pb.User
	}{
		{
			Name:     "MissingEmail",
			Password: "passwordabc",
			Error:    handler.ErrMissingEmail,
		},
		{
			Name:  "MissingPassword",
			Email: "john@doe.com",
			Error: handler.ErrInvalidPassword,
		},
		{
			Name:     "UserNotFound",
			Email:    "foo@bar.com",
			Password: "passwordabc",
			Error:    handler.ErrNotFound,
		},
		{
			Name:     "IncorrectPassword",
			Email:    "john@doe.com",
			Password: "passwordabcdef",
			Error:    handler.ErrIncorrectPassword,
		},
		{
			Name:     "Valid",
			Email:    "john@doe.com",
			Password: "passwordabc",
			User:     cRsp.User,
		},
	}

	for _, tc := range tt {
		t.Run(tc.Name, func(t *testing.T) {
			var rsp pb.LoginResponse
			err := h.Login(microAccountCtx(), &pb.LoginRequest{
				Email: tc.Email, Password: tc.Password,
			}, &rsp)
			assert.Equal(t, tc.Error, err)

			if tc.User != nil {
				assertUsersMatch(t, tc.User, rsp.User)
				assert.NotEmpty(t, rsp.Token)
			} else {
				assert.Nil(t, tc.User)
			}
		})
	}
}
