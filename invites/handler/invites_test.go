package handler_test

import (
	"context"
	"testing"

	"google.golang.org/protobuf/types/known/wrapperspb"

	"github.com/micro/services/invites/handler"
	pb "github.com/micro/services/invites/proto"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func testHandler(t *testing.T) *handler.Invites {
	// connect to the database
	db, err := gorm.Open(postgres.Open("postgresql://postgres@localhost:5432/postgres?sslmode=disable"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Error connecting to database: %v", err)
	}

	// clean any data from a previous run
	if err := db.Exec("DROP TABLE IF EXISTS invites CASCADE").Error; err != nil {
		t.Fatalf("Error cleaning database: %v", err)
	}

	// migrate the database
	if err := db.AutoMigrate(&handler.Invite{}); err != nil {
		t.Fatalf("Error migrating database: %v", err)
	}

	return &handler.Invites{DB: db}
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
			err := h.Create(context.TODO(), &pb.CreateRequest{
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
	err := h.Create(context.TODO(), &pb.CreateRequest{Email: "john@doe.com", GroupId: uuid.New().String()}, &cRsp)
	assert.NoError(t, err)
	if cRsp.Invite == nil {
		t.Fatal("No invite returned on create")
		return
	}

	tt := []struct {
		Name   string
		ID     *wrapperspb.StringValue
		Code   *wrapperspb.StringValue
		Error  error
		Invite *pb.Invite
	}{
		{
			Name:  "MissingIDAndCode",
			Error: handler.ErrMissingIDAndCode,
		},
		{
			Name:  "NotFoundByID",
			ID:    &wrapperspb.StringValue{Value: uuid.New().String()},
			Error: handler.ErrInviteNotFound,
		},
		{
			Name:  "NotFoundByCode",
			Code:  &wrapperspb.StringValue{Value: "12345678"},
			Error: handler.ErrInviteNotFound,
		},
		{
			Name:   "ValidID",
			ID:     &wrapperspb.StringValue{Value: cRsp.Invite.Id},
			Invite: cRsp.Invite,
		},
		{
			Name:   "ValidCode",
			Code:   &wrapperspb.StringValue{Value: cRsp.Invite.Code},
			Invite: cRsp.Invite,
		},
	}

	for _, tc := range tt {
		t.Run(tc.Name, func(t *testing.T) {
			var rsp pb.ReadResponse
			err := h.Read(context.TODO(), &pb.ReadRequest{Id: tc.ID, Code: tc.Code}, &rsp)
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
	err := h.Create(context.TODO(), &pb.CreateRequest{Email: "john@doe.com", GroupId: uuid.New().String()}, &cRsp)
	assert.NoError(t, err)
	if cRsp.Invite == nil {
		t.Fatal("No invite returned on create")
		return
	}

	tt := []struct {
		Name    string
		GroupID *wrapperspb.StringValue
		Email   *wrapperspb.StringValue
		Error   error
		Invite  *pb.Invite
	}{
		{
			Name:  "MissingIDAndEmail",
			Error: handler.ErrMissingGroupIDAndEmail,
		},
		{
			Name:  "NoResultsForEmail",
			Email: &wrapperspb.StringValue{Value: "foo@bar.com"},
		},
		{
			Name:    "NoResultsForGroupID",
			GroupID: &wrapperspb.StringValue{Value: uuid.New().String()},
		},
		{
			Name:    "ValidGroupID",
			GroupID: &wrapperspb.StringValue{Value: cRsp.Invite.GroupId},
			Invite:  cRsp.Invite,
		},
		{
			Name:   "ValidEmail",
			Email:  &wrapperspb.StringValue{Value: cRsp.Invite.Email},
			Invite: cRsp.Invite,
		},
		{
			Name:    "EmailAndGroupID",
			Email:   &wrapperspb.StringValue{Value: cRsp.Invite.Email},
			GroupID: &wrapperspb.StringValue{Value: uuid.New().String()},
		},
	}

	for _, tc := range tt {
		t.Run(tc.Name, func(t *testing.T) {
			var rsp pb.ListResponse
			err := h.List(context.TODO(), &pb.ListRequest{Email: tc.Email, GroupId: tc.GroupID}, &rsp)
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
		err := h.Delete(context.TODO(), &pb.DeleteRequest{}, &pb.DeleteResponse{})
		assert.Equal(t, handler.ErrMissingID, err)
	})

	// seed some data
	var cRsp pb.CreateResponse
	err := h.Create(context.TODO(), &pb.CreateRequest{Email: "john@doe.com", GroupId: uuid.New().String()}, &cRsp)
	assert.NoError(t, err)
	if cRsp.Invite == nil {
		t.Fatal("No invite returned on create")
		return
	}

	t.Run("Valid", func(t *testing.T) {
		err := h.Delete(context.TODO(), &pb.DeleteRequest{Id: cRsp.Invite.Id}, &pb.DeleteResponse{})
		assert.NoError(t, err)

		err = h.Read(context.TODO(), &pb.ReadRequest{Id: &wrapperspb.StringValue{Value: cRsp.Invite.Id}}, &pb.ReadResponse{})
		assert.Equal(t, handler.ErrInviteNotFound, err)
	})

	t.Run("Repeat", func(t *testing.T) {
		err := h.Delete(context.TODO(), &pb.DeleteRequest{Id: cRsp.Invite.Id}, &pb.DeleteResponse{})
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
