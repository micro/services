package handler

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/micro/micro/v3/service/auth"
	"github.com/micro/micro/v3/service/errors"
	"github.com/micro/micro/v3/service/store"
	pb "github.com/micro/services/groups/proto"
	"github.com/micro/services/pkg/tenant"
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
	ID      string
	Name    string
	Members []string
}

type Member struct {
	ID    string
	Group string
}

func (g *Group) Key(ctx context.Context) string {
	key := fmt.Sprintf("group:%s", g.ID)

	t, ok := tenant.FromContext(ctx)
	if !ok {
		return key
	}

	return fmt.Sprintf("%s/%s", t, key)
}

func (m *Member) Key(ctx context.Context) string {
	key := fmt.Sprintf("member:%s:%s", m.ID, m.Group)

	t, ok := tenant.FromContext(ctx)
	if !ok {
		return key
	}

	return fmt.Sprintf("%s/%s", t, key)

}

func (g *Group) Serialize() *pb.Group {
	memberIDs := make([]string, len(g.Members))
	for i, m := range g.Members {
		memberIDs[i] = m
	}
	return &pb.Group{Id: g.ID, Name: g.Name, MemberIds: memberIDs}
}

type Groups struct{}

func (g *Groups) Create(ctx context.Context, req *pb.CreateRequest, rsp *pb.CreateResponse) error {
	_, ok := auth.AccountFromContext(ctx)
	if !ok {
		errors.Unauthorized("UNAUTHORIZED", "Unauthorized")
	}
	// validate the request
	if len(req.Name) == 0 {
		return ErrMissingName
	}

	// create the group object
	group := &Group{ID: uuid.New().String(), Name: req.Name}

	// write the group record
	if err := store.Write(store.NewRecord(group.Key(ctx), group)); err != nil {
		return ErrStore
	}

	// return the group
	rsp.Group = group.Serialize()

	return nil
}

func (g *Groups) Read(ctx context.Context, req *pb.ReadRequest, rsp *pb.ReadResponse) error {
	_, ok := auth.AccountFromContext(ctx)
	if !ok {
		errors.Unauthorized("UNAUTHORIZED", "Unauthorized")
	}
	// validate the request
	if len(req.Ids) == 0 {
		return ErrMissingIDs
	}

	// serialize the response
	rsp.Groups = make(map[string]*pb.Group)

	for _, id := range req.Ids {
		group := &Group{
			ID: id,
		}
		recs, err := store.Read(group.Key(ctx), store.ReadLimit(1))
		if err != nil {
			return ErrStore
		}
		if len(recs) == 0 {
			continue
		}
		if err := recs[0].Decode(&group); err != nil {
			continue
		}
		rsp.Groups[group.ID] = group.Serialize()
	}

	return nil
}

func (g *Groups) Update(ctx context.Context, req *pb.UpdateRequest, rsp *pb.UpdateResponse) error {
	_, ok := auth.AccountFromContext(ctx)
	if !ok {
		errors.Unauthorized("UNAUTHORIZED", "Unauthorized")
	}
	// validate the request
	if len(req.Id) == 0 {
		return ErrMissingID
	}
	if len(req.Name) == 0 {
		return ErrMissingName
	}

	group := &Group{ID: req.Id}

	recs, err := store.Read(group.Key(ctx), store.ReadLimit(1))
	if err == store.ErrNotFound {
		return ErrNotFound
	} else if err != nil {
		return ErrStore
	}

	// decode the record
	recs[0].Decode(&group)

	// set the name
	group.Name = req.Name

	// save the record
	if err := store.Write(store.NewRecord(group.Key(ctx), group)); err != nil {
		return ErrStore
	}

	// serialize the response
	rsp.Group = group.Serialize()

	return nil
}

func (g *Groups) Delete(ctx context.Context, req *pb.DeleteRequest, rsp *pb.DeleteResponse) error {
	_, ok := auth.AccountFromContext(ctx)
	if !ok {
		errors.Unauthorized("UNAUTHORIZED", "Unauthorized")
	}
	// validate the request
	if len(req.Id) == 0 {
		return ErrMissingID
	}

	group := &Group{ID: req.Id}

	// get the group
	recs, err := store.Read(group.Key(ctx), store.ReadLimit(1))
	if err == store.ErrNotFound {
		return nil
	} else if err != nil {
		return ErrStore
	}

	// decode the record
	recs[0].Decode(&group)

	// delete the record
	if err := store.Delete(group.Key(ctx)); err == store.ErrNotFound {
		return nil
	} else if err != nil {
		return ErrStore
	}

	// delete all the members
	for _, memberId := range group.Members {
		m := &Member{
			ID: memberId,
		}
		// delete the member
		store.Delete(m.Key(ctx))
	}

	return nil
}

func (g *Groups) List(ctx context.Context, req *pb.ListRequest, rsp *pb.ListResponse) error {
	_, ok := auth.AccountFromContext(ctx)
	if !ok {
		errors.Unauthorized("UNAUTHORIZED", "Unauthorized")
	}

	if len(req.MemberId) > 0 {
		// only list groups the user is a member of
		m := &Member{ID: req.MemberId}
		recs, err := store.Read(m.Key(ctx), store.ReadPrefix())
		if err != nil {
			return ErrStore
		}

		for _, rec := range recs {
			m := &Member{ID: req.MemberId}
			rec.Decode(&m)

			// get the group
			group := &Group{ID: m.Group}

			grecs, err := store.Read(group.Key(ctx), store.ReadLimit(1))
			if err != nil {
				return ErrStore
			}
			grecs[0].Decode(&group)

			rsp.Groups = append(rsp.Groups, group.Serialize())
		}

		return nil
	}

	group := &Group{}

	// read all the prefixes
	recs, err := store.Read(group.Key(ctx), store.ReadPrefix())
	if err != nil {
		return ErrStore
	}

	// serialize and return response
	for _, rec := range recs {
		group := new(Group)
		rec.Decode(&group)
		rsp.Groups = append(rsp.Groups, group.Serialize())
	}

	return nil
}

func (g *Groups) AddMember(ctx context.Context, req *pb.AddMemberRequest, rsp *pb.AddMemberResponse) error {
	_, ok := auth.AccountFromContext(ctx)
	if !ok {
		errors.Unauthorized("UNAUTHORIZED", "Unauthorized")
	}
	// validate the request
	if len(req.GroupId) == 0 {
		return ErrMissingGroupID
	}
	if len(req.MemberId) == 0 {
		return ErrMissingMemberID
	}

	// read the group
	group := &Group{ID: req.GroupId}

	recs, err := store.Read(group.Key(ctx), store.ReadLimit(1))
	if err == store.ErrNotFound {
		return ErrNotFound
	} else if err != nil {
		return ErrStore
	}

	// decode the record
	recs[0].Decode(group)

	var seen bool
	for _, member := range group.Members {
		if member == req.MemberId {
			seen = true
			break
		}
	}

	// already a member
	if seen {
		return nil
	}

	// add the member
	group.Members = append(group.Members, req.MemberId)

	// save the record
	if err := store.Write(store.NewRecord(group.Key(ctx), group)); err != nil {
		return ErrStore
	}

	// add the member record

	m := &Member{
		ID:    req.MemberId,
		Group: group.ID,
	}

	// write the record
	if err := store.Write(store.NewRecord(m.Key(ctx), m)); err != nil {
		return ErrStore
	}

	return nil
}

func (g *Groups) RemoveMember(ctx context.Context, req *pb.RemoveMemberRequest, rsp *pb.RemoveMemberResponse) error {
	_, ok := auth.AccountFromContext(ctx)
	if !ok {
		errors.Unauthorized("UNAUTHORIZED", "Unauthorized")
	}
	// validate the request
	if len(req.GroupId) == 0 {
		return ErrMissingGroupID
	}
	if len(req.MemberId) == 0 {
		return ErrMissingMemberID
	}

	// read the group
	group := &Group{ID: req.GroupId}

	// read the gruop
	recs, err := store.Read(group.Key(ctx), store.ReadLimit(1))
	if err == store.ErrNotFound {
		return ErrNotFound
	} else if err != nil {
		return ErrStore
	}

	// decode the record
	recs[0].Decode(&group)

	// new member id list
	var members []string

	for _, member := range group.Members {
		if member == req.MemberId {
			continue
		}
		members = append(members, member)
	}

	// update the member
	group.Members = members

	// save the record
	if err := store.Write(store.NewRecord(group.Key(ctx), group)); err != nil {
		return ErrStore
	}

	// delete the member
	m := &Member{
		ID:    req.MemberId,
		Group: group.ID,
	}
	if err := store.Delete(m.Key(ctx)); err != nil {
		return ErrStore
	}

	return nil
}
