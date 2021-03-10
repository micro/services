package handler_test

import (
	"context"
	"testing"

	"github.com/micro/services/users/handler"
	pb "github.com/micro/services/users/proto"
	"github.com/stretchr/testify/assert"
)

func TestCreate(t *testing.T) {
	tt := []struct {
		Name      string
		FirstName string
		LastName  string
		Email     string
		Password  string
		Error     error
	}{
		{
			Name:     "MissingFirstName",
			LastName: "Doe",
			Email:    "john@doe.com",
			Password: "password",
			Error:    handler.ErrMissingFirstName,
		},
		{
			Name:      "MissingLastName",
			FirstName: "John",
			Email:     "john@doe.com",
			Password:  "password",
			Error:     handler.ErrMissingLastName,
		},
		{
			Name:      "MissingEmail",
			FirstName: "John",
			LastName:  "Doe",
			Password:  "password",
			Error:     handler.ErrMissingEmail,
		},
		{
			Name:      "InvalidEmail",
			FirstName: "John",
			LastName:  "Doe",
			Password:  "password",
			Email:     "foo.foo.foo",
			Error:     handler.ErrInvalidEmail,
		},
		{
			Name:      "InvalidPassword",
			FirstName: "John",
			LastName:  "Doe",
			Email:     "john@doe.com",
			Password:  "pwd",
			Error:     handler.ErrInvalidPassword,
		},
	}

	// test the validations
	h := testHandler(t)
	for _, tc := range tt {
		t.Run(tc.Name, func(t *testing.T) {
			err := h.Create(context.TODO(), &pb.CreateRequest{
				FirstName: tc.FirstName,
				LastName:  tc.LastName,
				Email:     tc.Email,
				Password:  tc.Password,
			}, &pb.CreateResponse{})
			assert.Equal(t, tc.Error, err)
		})
	}

	t.Run("Valid", func(t *testing.T) {
		var rsp pb.CreateResponse
		req := pb.CreateRequest{
			FirstName: "John",
			LastName:  "Doe",
			Email:     "john@doe.com",
			Password:  "passwordabc",
		}
		err := h.Create(context.TODO(), &req, &rsp)

		assert.NoError(t, err)
		u := rsp.User
		if u == nil {
			t.Fatalf("No user returned")
		}
		assert.NotEmpty(t, u.Id)
		assert.Equal(t, req.FirstName, u.FirstName)
		assert.Equal(t, req.LastName, u.LastName)
		assert.Equal(t, req.Email, u.Email)
		assert.NotEmpty(t, rsp.Token)
	})

	t.Run("DuplicateEmail", func(t *testing.T) {
		var rsp pb.CreateResponse
		req := pb.CreateRequest{
			FirstName: "John",
			LastName:  "Doe",
			Email:     "john@doe.com",
			Password:  "passwordabc",
		}
		err := h.Create(context.TODO(), &req, &rsp)
		assert.Equal(t, handler.ErrDuplicateEmail, err)
		assert.Nil(t, rsp.User)
	})

	t.Run("DifferentEmail", func(t *testing.T) {
		var rsp pb.CreateResponse
		req := pb.CreateRequest{
			FirstName: "John",
			LastName:  "Doe",
			Email:     "johndoe@gmail.com",
			Password:  "passwordabc",
		}
		err := h.Create(context.TODO(), &req, &rsp)

		assert.NoError(t, err)
		u := rsp.User
		if u == nil {
			t.Fatalf("No user returned")
		}
		assert.NotEmpty(t, u.Id)
		assert.Equal(t, req.FirstName, u.FirstName)
		assert.Equal(t, req.LastName, u.LastName)
		assert.Equal(t, req.Email, u.Email)
	})
}
