package handler

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"github.com/micro/micro/v3/service/errors"
	pb "github.com/micro/services/groups/proto"
	"gorm.io/gorm"
)

var (
	ErrMissingName     = errors.BadRequest("MISSING_NAME", "Missing name")
	ErrMissingID       = errors.BadRequest("MISSING_ID", "Missing ID")
	ErrMissingIDs      = errors.BadRequest("MISSING_IDS", "One or more IDs are required")
	ErrMissingGroupID  = errors.BadRequest("MISSING_GROUP_ID", "Missing Group ID")
	ErrMissingMemberID = errors.BadRequest("MISSING_MEMBER_ID", "Missing Member ID")
	ErrNotFound        = errors.BadRequest("NOT_FOUND", "No group found with this ID")
	ErrStore           = errors.InternalServerError("STORE_ERROR", "Error connecting to the store")
)

type Group struct {
	ID          string
	Name        string
	Memberships []Membership
}

type Membership struct {
	MemberID string `gorm:"uniqueIndex:idx_membership"`
	GroupID  string `gorm:"uniqueIndex:idx_membership"`
	Group    Group
}

func (g *Group) Serialize() *pb.Group {
	memberIDs := make([]string, len(g.Memberships))
	for i, m := range g.Memberships {
		memberIDs[i] = m.MemberID
	}
	return &pb.Group{Id: g.ID, Name: g.Name, MemberIds: memberIDs}
}

type Groups struct {
	DB *gorm.DB
}

func (g *Groups) Create(ctx context.Context, req *pb.CreateRequest, rsp *pb.CreateResponse) error {
	// validate the request
	if len(req.Name) == 0 {
		return ErrMissingName
	}

	// create the group object
	group := &Group{ID: uuid.New().String(), Name: req.Name}
	if err := g.DB.Create(group).Error; err != nil {
		return ErrStore
	}

	// return the group
	rsp.Group = group.Serialize()
	return nil
}

func (g *Groups) Read(ctx context.Context, req *pb.ReadRequest, rsp *pb.ReadResponse) error {
	// validate the request
	if len(req.Ids) == 0 {
		return ErrMissingIDs
	}

	// query the database
	var groups []Group
	if err := g.DB.Model(&Group{}).Preload("Memberships").Where("id IN (?)", req.Ids).Find(&groups).Error; err != nil {
		return ErrStore
	}

	// serialize the response
	rsp.Groups = make(map[string]*pb.Group, len(groups))
	for _, g := range groups {
		rsp.Groups[g.ID] = g.Serialize()
	}

	return nil
}

func (g *Groups) Update(ctx context.Context, req *pb.UpdateRequest, rsp *pb.UpdateResponse) error {
	// validate the request
	if len(req.Id) == 0 {
		return ErrMissingID
	}
	if len(req.Name) == 0 {
		return ErrMissingName
	}

	return g.DB.Transaction(func(tx *gorm.DB) error {
		// find the group
		var group Group
		if err := tx.Where(&Group{ID: req.Id}).First(&group).Error; err == gorm.ErrRecordNotFound {
			return ErrNotFound
		} else if err != nil {
			return ErrStore
		}

		// update the group
		group.Name = req.Name
		if err := tx.Save(&group).Error; err != nil {
			return ErrStore
		}

		// serialize the response
		rsp.Group = group.Serialize()
		return nil
	})
}

func (g *Groups) Delete(ctx context.Context, req *pb.DeleteRequest, rsp *pb.DeleteResponse) error {
	// validate the request
	if len(req.Id) == 0 {
		return ErrMissingID
	}

	// delete from the database
	if err := g.DB.Delete(&Group{ID: req.Id}).Error; err == gorm.ErrRecordNotFound {
		return nil
	} else if err != nil {
		return ErrStore
	}

	return nil
}

func (g *Groups) List(ctx context.Context, req *pb.ListRequest, rsp *pb.ListResponse) error {
	if len(req.MemberId) > 0 {
		// only list groups the user is a member of
		var ms []Membership
		q := g.DB.Where(&Membership{MemberID: req.MemberId}).Preload("Group.Memberships")
		if err := q.Find(&ms).Error; err != nil {
			return err
		}
		rsp.Groups = make([]*pb.Group, len(ms))
		for i, m := range ms {
			rsp.Groups[i] = m.Group.Serialize()
		}
		return nil
	}

	// load all groups
	var groups []Group
	if err := g.DB.Model(&Group{}).Preload("Memberships").Find(&groups).Error; err != nil {
		return ErrStore
	}

	// serialize the response
	rsp.Groups = make([]*pb.Group, len(groups))
	for i, g := range groups {
		rsp.Groups[i] = g.Serialize()
	}

	return nil
}

func (g *Groups) AddMember(ctx context.Context, req *pb.AddMemberRequest, rsp *pb.AddMemberResponse) error {
	// validate the request
	if len(req.GroupId) == 0 {
		return ErrMissingGroupID
	}
	if len(req.MemberId) == 0 {
		return ErrMissingMemberID
	}

	return g.DB.Transaction(func(tx *gorm.DB) error {
		// check the group exists
		var group Group
		if err := tx.Where(&Group{ID: req.GroupId}).First(&group).Error; err == gorm.ErrRecordNotFound {
			return ErrNotFound
		} else if err != nil {
			return err
		}

		// create the membership
		m := &Membership{MemberID: req.MemberId, GroupID: req.GroupId}
		err := tx.Create(m).Error
		// check for membership already existing (unique index violation)
		if err != nil && strings.Contains(err.Error(), "fk_groups_memberships") {
			return nil
		} else if err != nil {
			return ErrStore
		}

		return nil
	})
}

func (g *Groups) RemoveMember(ctx context.Context, req *pb.RemoveMemberRequest, rsp *pb.RemoveMemberResponse) error {
	// validate the request
	if len(req.GroupId) == 0 {
		return ErrMissingGroupID
	}
	if len(req.MemberId) == 0 {
		return ErrMissingMemberID
	}

	// delete the membership
	m := &Membership{MemberID: req.MemberId, GroupID: req.GroupId}
	if err := g.DB.Where(m).Delete(m).Error; err != nil {
		return ErrStore
	}

	return nil
}
