package handler_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/micro/services/users/handler"
	pb "github.com/micro/services/users/proto"
	"github.com/stretchr/testify/assert"
)

func TestValidate(t *testing.T) {
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
		FirstName: "Barry",
		LastName:  "Doe",
		Email:     "barry@doe.com",
		Password:  "passwordabc",
	}
	err = h.Create(microAccountCtx(), &cReq2, &cRsp2)
	assert.NoError(t, err)
	if cRsp2.User == nil {
		t.Fatal("No user returned")
		return
	}

	tt := []struct {
		Name  string
		Token string
		Time  func() time.Time
		Error error
		User  *pb.User
	}{
		{
			Name:  "MissingToken",
			Error: handler.ErrMissingToken,
		},
		{
			Name:  "InvalidToken",
			Error: handler.ErrInvalidToken,
			Token: uuid.New().String(),
		},
		{
			Name:  "ExpiredToken",
			Error: handler.ErrTokenExpired,
			Token: cRsp1.Token,
			Time:  func() time.Time { return time.Now().Add(time.Hour * 24 * 8) },
		},
		{
			Name:  "ValidToken",
			User:  cRsp2.User,
			Token: cRsp2.Token,
			Time:  func() time.Time { return time.Now().Add(time.Hour * 24 * 3) },
		},
		{
			Name:  "RefreshedToken",
			User:  cRsp2.User,
			Token: cRsp2.Token,
			Time:  func() time.Time { return time.Now().Add(time.Hour * 24 * 8) },
		},
	}

	for _, tc := range tt {
		t.Run(tc.Name, func(t *testing.T) {
			if tc.Time == nil {
				h.Time = time.Now
			} else {
				h.Time = tc.Time
			}

			var rsp pb.ValidateResponse
			err := h.Validate(microAccountCtx(), &pb.ValidateRequest{Token: tc.Token}, &rsp)
			assert.Equal(t, tc.Error, err)

			if tc.User != nil {
				assertUsersMatch(t, tc.User, rsp.User)
			} else {
				assert.Nil(t, tc.User)
			}
		})
	}
}
