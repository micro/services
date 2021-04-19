package handler_test

import (
	"context"
	"database/sql"
	"os"
	"testing"

	"github.com/micro/micro/v3/service/auth"
	"github.com/micro/services/invites/handler"
	pb "github.com/micro/services/invites/proto"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func testHandler(t *testing.T) *handler.Invites {
	// connect to the database
	addr := os.Getenv("POSTGRES_URL")
	if len(addr) == 0 {
		addr = "postgresql://postgres@localhost:5432/postgres?sslmode=disable"
	}
	sqlDB, err := sql.Open("pgx", addr)
	if err != nil {
		t.Fatalf("Failed to open connection to DB %s", err)
	}

	// clean any data from a previous run
	if _, err := sqlDB.Exec("DROP TABLE IF EXISTS micro_invites CASCADE"); err != nil {
		t.Fatalf("Error cleaning database: %v", err)
	}

	h := &handler.Invites{}
	h.DBConn(sqlDB).Migrations(&handler.Invite{})
	return h
}

func TestCreate(t *testing.T) {
	tt := []struct {
		Name    string
		GroupID string
		Email   string
		Error   error
	}{
		{
			Name:  "MissingGroupID",
			Email: "john@doe.com",
			Error: handler.ErrMissingGroupID,
		},
		{
			Name:    "MissingEmail",
			GroupID: uuid.New().String(),
			Error:   handler.ErrMissingEmail,
		},
		{
			Name:    "InvalidEmail",
			GroupID: uuid.New().String(),
			Email:   "foo.foo.foo",
			Error:   handler.ErrInvalidEmail,
		},
		{
			Name:    "Valid",
			GroupID: "thisisavalidgroupid",
			Email:   "john@doe.com",
		},
		{
			Name:    "Repeat",
			GroupID: "thisisavalidgroupid",
			Email:   "john@doe.com",
		},
	}

	h := testHandler(t)
	for _, tc := range tt {
		t.Run(tc.Name, func(t *testing.T) {
			var rsp pb.CreateResponse
			err := h.Create(microAccountCtx(), &pb.CreateRequest{
				GroupId: tc.GroupID, Email: tc.Email,
			}, &rsp)
			assert.Equal(t, tc.Error, err)

			if tc.Error != nil {
				assert.Nil(t, rsp.Invite)
				return
			}

			if rsp.Invite == nil {
				t.Fatalf("Invite was not returned")
				return
			}

			assert.NotEmpty(t, rsp.Invite.Id)
			assert.NotEmpty(t, rsp.Invite.Code)
			assert.Equal(t, tc.GroupID, rsp.Invite.GroupId)
			assert.Equal(t, tc.Email, rsp.Invite.Email)
		})
	}
}

func TestRead(t *testing.T) {
	h := testHandler(t)

	// seed some data
	var cRsp pb.CreateResponse
	err := h.Create(microAccountCtx(), &pb.CreateRequest{Email: "john@doe.com", GroupId: uuid.New().String()}, &cRsp)
	assert.NoError(t, err)
	if cRsp.Invite == nil {
		t.Fatal("No invite returned on create")
		return
	}

	tt := []struct {
		Name   string
		ID     string
		Code   string
		Error  error
		Invite *pb.Invite
	}{
		{
			Name:  "MissingIDAndCode",
			Error: handler.ErrMissingIDAndCode,
		},
		{
			Name:  "NotFoundByID",
			ID:    uuid.New().String(),
			Error: handler.ErrInviteNotFound,
		},
		{
			Name:  "NotFoundByCode",
			Code:  "12345678",
			Error: handler.ErrInviteNotFound,
		},
		{
			Name:   "ValidID",
			ID:     cRsp.Invite.Id,
			Invite: cRsp.Invite,
		},
		{
			Name:   "ValidCode",
			Code:   cRsp.Invite.Code,
			Invite: cRsp.Invite,
		},
	}

	for _, tc := range tt {
		t.Run(tc.Name, func(t *testing.T) {
			var rsp pb.ReadResponse
			err := h.Read(microAccountCtx(), &pb.ReadRequest{Id: tc.ID, Code: tc.Code}, &rsp)
			assert.Equal(t, tc.Error, err)

			if tc.Invite == nil {
				assert.Nil(t, rsp.Invite)
			} else {
				assertInvitesMatch(t, tc.Invite, rsp.Invite)
			}
		})
	}
}

func TestList(t *testing.T) {
	h := testHandler(t)

	// seed some data
	var cRsp pb.CreateResponse
	err := h.Create(microAccountCtx(), &pb.CreateRequest{Email: "john@doe.com", GroupId: uuid.New().String()}, &cRsp)
	assert.NoError(t, err)
	if cRsp.Invite == nil {
		t.Fatal("No invite returned on create")
		return
	}

	tt := []struct {
		Name    string
		GroupID string
		Email   string
		Error   error
		Invite  *pb.Invite
	}{
		{
			Name:  "MissingIDAndEmail",
			Error: handler.ErrMissingGroupIDAndEmail,
		},
		{
			Name:  "NoResultsForEmail",
			Email: "foo@bar.com",
		},
		{
			Name:    "NoResultsForGroupID",
			GroupID: uuid.New().String(),
		},
		{
			Name:    "ValidGroupID",
			GroupID: cRsp.Invite.GroupId,
			Invite:  cRsp.Invite,
		},
		{
			Name:   "ValidEmail",
			Email:  cRsp.Invite.Email,
			Invite: cRsp.Invite,
		},
		{
			Name:    "EmailAndGroupID",
			Email:   cRsp.Invite.Email,
			GroupID: uuid.New().String(),
		},
	}

	for _, tc := range tt {
		t.Run(tc.Name, func(t *testing.T) {
			var rsp pb.ListResponse
			err := h.List(microAccountCtx(), &pb.ListRequest{Email: tc.Email, GroupId: tc.GroupID}, &rsp)
			assert.Equal(t, tc.Error, err)

			if tc.Invite == nil {
				assert.Empty(t, rsp.Invites)
			} else {
				if len(rsp.Invites) != 1 {
					t.Errorf("Incorrect number of invites returned, expected 1 but got %v", len(rsp.Invites))
					return
				}
				assertInvitesMatch(t, tc.Invite, rsp.Invites[0])
			}
		})
	}
}

func TestDelete(t *testing.T) {
	h := testHandler(t)

	t.Run("MissingID", func(t *testing.T) {
		err := h.Delete(microAccountCtx(), &pb.DeleteRequest{}, &pb.DeleteResponse{})
		assert.Equal(t, handler.ErrMissingID, err)
	})

	// seed some data
	var cRsp pb.CreateResponse
	err := h.Create(microAccountCtx(), &pb.CreateRequest{Email: "john@doe.com", GroupId: uuid.New().String()}, &cRsp)
	assert.NoError(t, err)
	if cRsp.Invite == nil {
		t.Fatal("No invite returned on create")
		return
	}

	t.Run("Valid", func(t *testing.T) {
		err := h.Delete(microAccountCtx(), &pb.DeleteRequest{Id: cRsp.Invite.Id}, &pb.DeleteResponse{})
		assert.NoError(t, err)

		err = h.Read(microAccountCtx(), &pb.ReadRequest{Id: cRsp.Invite.Id}, &pb.ReadResponse{})
		assert.Equal(t, handler.ErrInviteNotFound, err)
	})

	t.Run("Repeat", func(t *testing.T) {
		err := h.Delete(microAccountCtx(), &pb.DeleteRequest{Id: cRsp.Invite.Id}, &pb.DeleteResponse{})
		assert.NoError(t, err)
	})
}

func assertInvitesMatch(t *testing.T, exp, act *pb.Invite) {
	if act == nil {
		t.Errorf("No invite returned")
		return
	}
	assert.Equal(t, exp.Id, act.Id)
	assert.Equal(t, exp.Code, act.Code)
	assert.Equal(t, exp.Email, act.Email)
	assert.Equal(t, exp.GroupId, act.GroupId)
}

func microAccountCtx() context.Context {
	return auth.ContextWithAccount(context.TODO(), &auth.Account{
		Issuer: "micro",
	})
}
