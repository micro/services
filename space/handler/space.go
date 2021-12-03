package handler

import (
	"context"
	"sync"
	"time"


	"github.com/micro/services/pkg/tenant"
	"github.com/micro/micro/v3/service/store"
	pb "github.com/micro/services/space/proto"
)

type Space struct{}

var (
	mtx sync.RWMutex

	voteKey = "votes/"
)

type Vote struct {
	Id string `json:"id"`
	Message string `json:"message"`
	VotedAt time.Time `json:"voted_at"`
}

func (n *Space) Vote(ctx context.Context, req *pb.VoteRequest, rsp *pb.VoteResponse) error {
	mtx.Lock()
	defer mtx.Unlock()

	id, ok := tenant.FromContext(ctx)
	if !ok {
		id = "micro"
	}

	rec := store.NewRecord(voteKey + id, &Vote{
		Id: id,
		Message: req.Message,
		VotedAt: time.Now(),
	})

	// we don't need to check the error
	store.Write(rec)

	rsp.Message = "Thanks for the vote!"

	return nil
}

