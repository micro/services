package handler_test

import (
	"context"
	"os"
	"sort"
	"testing"

	"github.com/google/uuid"
	"github.com/micro/services/groups/handler"
	pb "github.com/micro/services/groups/proto"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func testHandler(t *testing.T) *handler.Groups {
	// connect to the database
	addr := os.Getenv("POSTGRES_URL")
	if len(addr) == 0 {
		addr = "postgresql://postgres@localhost:5432/postgres?sslmode=disable"
	}
	db, err := gorm.Open(postgres.Open(addr), &gorm.Config{})
	if err != nil {
		t.Fatalf("Error connecting to database: %v", err)
	}

	// clean any data from a previous run
	if err := db.Exec("DROP TABLE IF EXISTS groups, memberships CASCADE").Error; err != nil {
		t.Fatalf("Error cleaning database: %v", err)
	}

	// migrate the database
	if err := db.AutoMigrate(&handler.Group{}, &handler.Membership{}); err != nil {
		t.Fatalf("Error migrating database: %v", err)
	}

	return &handler.Groups{DB: db}
}
func TestCreate(t *testing.T) {
	h := testHandler(t)

	t.Run("MissingName", func(t *testing.T) {
		err := h.Create(context.TODO(), &pb.CreateRequest{}, &pb.CreateResponse{})
		assert.Equal(t, handler.ErrMissingName, err)
	})

	t.Run("Valid", func(t *testing.T) {
		err := h.Create(context.TODO(), &pb.CreateRequest{
			Name: "Doe Family Group",
		}, &pb.CreateResponse{})
		assert.NoError(t, err)
	})
}

func TestUpdate(t *testing.T) {
	h := testHandler(t)

	t.Run("MissingID", func(t *testing.T) {
		err := h.Update(context.TODO(), &pb.UpdateRequest{
			Name: "Doe Family Group",
		}, &pb.UpdateResponse{})
		assert.Equal(t, handler.ErrMissingID, err)
	})

	t.Run("MissingName", func(t *testing.T) {
		err := h.Update(context.TODO(), &pb.UpdateRequest{
			Id: uuid.New().String(),
		}, &pb.UpdateResponse{})
		assert.Equal(t, handler.ErrMissingName, err)
	})

	t.Run("NotFound", func(t *testing.T) {
		err := h.Update(context.TODO(), &pb.UpdateRequest{
			Id:   uuid.New().String(),
			Name: "Bar Family Group",
		}, &pb.UpdateResponse{})
		assert.Equal(t, handler.ErrNotFound, err)
	})

	t.Run("Valid", func(t *testing.T) {
		// create a demo group
		var cRsp pb.CreateResponse
		err := h.Create(context.TODO(), &pb.CreateRequest{
			Name: "Doe Family Group",
		}, &cRsp)
		assert.NoError(t, err)

		err = h.Update(context.TODO(), &pb.UpdateRequest{
			Id:   cRsp.Group.Id,
			Name: "Bar Family Group",
		}, &pb.UpdateResponse{})
		assert.NoError(t, err)

		var rRsp pb.ReadResponse
		err = h.Read(context.TODO(), &pb.ReadRequest{
			Ids: []string{cRsp.Group.Id},
		}, &rRsp)
		assert.NoError(t, err)

		g := rRsp.Groups[cRsp.Group.Id]
		if g == nil {
			t.Errorf("Group not returned")
		} else {
			assert.Equal(t, "Bar Family Group", g.Name)
		}
	})
}

func TestDelete(t *testing.T) {
	h := testHandler(t)

	t.Run("MissingID", func(t *testing.T) {
		err := h.Delete(context.TODO(), &pb.DeleteRequest{}, &pb.DeleteResponse{})
		assert.Equal(t, handler.ErrMissingID, err)
	})

	t.Run("NotFound", func(t *testing.T) {
		err := h.Delete(context.TODO(), &pb.DeleteRequest{
			Id: uuid.New().String(),
		}, &pb.DeleteResponse{})
		assert.NoError(t, err)
	})

	// create a demo group
	var cRsp pb.CreateResponse
	err := h.Create(context.TODO(), &pb.CreateRequest{
		Name: "Doe Family Group",
	}, &cRsp)
	assert.NoError(t, err)

	t.Run("Valid", func(t *testing.T) {
		err := h.Delete(context.TODO(), &pb.DeleteRequest{
			Id: cRsp.Group.Id,
		}, &pb.DeleteResponse{})
		assert.NoError(t, err)

		var rRsp pb.ReadResponse
		err = h.Read(context.TODO(), &pb.ReadRequest{
			Ids: []string{cRsp.Group.Id},
		}, &rRsp)
		assert.Nil(t, rRsp.Groups[cRsp.Group.Id])
	})
}
func TestList(t *testing.T) {
	h := testHandler(t)

	// create two demo groups
	var cRsp1 pb.CreateResponse
	err := h.Create(context.TODO(), &pb.CreateRequest{
		Name: "Alpha Group",
	}, &cRsp1)
	assert.NoError(t, err)

	var cRsp2 pb.CreateResponse
	err = h.Create(context.TODO(), &pb.CreateRequest{
		Name: "Bravo Group",
	}, &cRsp2)
	assert.NoError(t, err)

	// add a member to the first group
	uid := uuid.New().String()
	err = h.AddMember(context.TODO(), &pb.AddMemberRequest{
		GroupId: cRsp1.Group.Id, MemberId: uid,
	}, &pb.AddMemberResponse{})
	assert.NoError(t, err)

	t.Run("Unscoped", func(t *testing.T) {
		var rsp pb.ListResponse
		err = h.List(context.TODO(), &pb.ListRequest{}, &rsp)
		assert.NoError(t, err)
		assert.Lenf(t, rsp.Groups, 2, "Two groups should be returned")
		if len(rsp.Groups) != 2 {
			return
		}

		sort.Slice(rsp.Groups, func(i, j int) bool {
			return rsp.Groups[i].Name < rsp.Groups[j].Name
		})
		assert.Equal(t, cRsp1.Group.Id, rsp.Groups[0].Id)
		assert.Equal(t, cRsp1.Group.Name, rsp.Groups[0].Name)
		assert.Len(t, rsp.Groups[0].MemberIds, 1)
		assert.Contains(t, rsp.Groups[0].MemberIds, uid)
		assert.Equal(t, cRsp2.Group.Id, rsp.Groups[1].Id)
		assert.Equal(t, cRsp2.Group.Name, rsp.Groups[1].Name)
		assert.Len(t, rsp.Groups[1].MemberIds, 0)
	})

	t.Run("Scoped", func(t *testing.T) {
		var rsp pb.ListResponse
		err = h.List(context.TODO(), &pb.ListRequest{MemberId: uid}, &rsp)
		assert.NoError(t, err)
		assert.Lenf(t, rsp.Groups, 1, "One group should be returned")
		if len(rsp.Groups) != 1 {
			return
		}

		assert.Equal(t, cRsp1.Group.Id, rsp.Groups[0].Id)
		assert.Equal(t, cRsp1.Group.Name, rsp.Groups[0].Name)
		assert.Len(t, rsp.Groups[0].MemberIds, 1)
		assert.Contains(t, rsp.Groups[0].MemberIds, uid)
	})
}

func TestAddMember(t *testing.T) {
	h := testHandler(t)

	t.Run("MissingGroupID", func(t *testing.T) {
		err := h.AddMember(context.TODO(), &pb.AddMemberRequest{
			MemberId: uuid.New().String(),
		}, &pb.AddMemberResponse{})
		assert.Equal(t, handler.ErrMissingGroupID, err)
	})

	t.Run("MissingMemberID", func(t *testing.T) {
		err := h.AddMember(context.TODO(), &pb.AddMemberRequest{
			GroupId: uuid.New().String(),
		}, &pb.AddMemberResponse{})
		assert.Equal(t, handler.ErrMissingMemberID, err)
	})

	t.Run("GroupNotFound", func(t *testing.T) {
		err := h.AddMember(context.TODO(), &pb.AddMemberRequest{
			GroupId:  uuid.New().String(),
			MemberId: uuid.New().String(),
		}, &pb.AddMemberResponse{})
		assert.Equal(t, handler.ErrNotFound, err)
	})

	// create a test group
	var cRsp pb.CreateResponse
	err := h.Create(context.TODO(), &pb.CreateRequest{
		Name: "Alpha Group",
	}, &cRsp)
	assert.NoError(t, err)

	t.Run("Valid", func(t *testing.T) {
		err := h.AddMember(context.TODO(), &pb.AddMemberRequest{
			GroupId:  cRsp.Group.Id,
			MemberId: uuid.New().String(),
		}, &pb.AddMemberResponse{})
		assert.NoError(t, err)
	})

	t.Run("Retry", func(t *testing.T) {
		err := h.AddMember(context.TODO(), &pb.AddMemberRequest{
			GroupId:  cRsp.Group.Id,
			MemberId: uuid.New().String(),
		}, &pb.AddMemberResponse{})
		assert.NoError(t, err)
	})
}

func TestRemoveMember(t *testing.T) {
	h := testHandler(t)

	t.Run("MissingGroupID", func(t *testing.T) {
		err := h.RemoveMember(context.TODO(), &pb.RemoveMemberRequest{
			MemberId: uuid.New().String(),
		}, &pb.RemoveMemberResponse{})
		assert.Equal(t, handler.ErrMissingGroupID, err)
	})

	t.Run("MissingMemberID", func(t *testing.T) {
		err := h.RemoveMember(context.TODO(), &pb.RemoveMemberRequest{
			GroupId: uuid.New().String(),
		}, &pb.RemoveMemberResponse{})
		assert.Equal(t, handler.ErrMissingMemberID, err)
	})

	// create a test group
	var cRsp pb.CreateResponse
	err := h.Create(context.TODO(), &pb.CreateRequest{
		Name: "Alpha Group",
	}, &cRsp)
	assert.NoError(t, err)

	t.Run("Valid", func(t *testing.T) {
		err := h.RemoveMember(context.TODO(), &pb.RemoveMemberRequest{
			GroupId:  cRsp.Group.Id,
			MemberId: uuid.New().String(),
		}, &pb.RemoveMemberResponse{})
		assert.NoError(t, err)
	})

	t.Run("Retry", func(t *testing.T) {
		err := h.RemoveMember(context.TODO(), &pb.RemoveMemberRequest{
			GroupId:  cRsp.Group.Id,
			MemberId: uuid.New().String(),
		}, &pb.RemoveMemberResponse{})
		assert.NoError(t, err)
	})
}
