package handler

import (
	"context"

	"github.com/micro/micro/v3/service/auth"
	"github.com/micro/micro/v3/service/errors"
	"github.com/micro/micro/v3/service/logger"
	"github.com/micro/services/pkg/model"
	pb "github.com/micro/services/threads/proto"
)

// List all the threads for a group
func (s *Threads) ListThreads(ctx context.Context, req *pb.ListThreadsRequest, rsp *pb.ListThreadsResponse) error {
	_, ok := auth.AccountFromContext(ctx)
	if !ok {
		errors.Unauthorized("UNAUTHORIZED", "Unauthorized")
	}
	// validate the request
	if len(req.GroupId) == 0 {
		return ErrMissingGroupID
	}

	var threads []*Thread
	thread := &Thread{GroupID: req.GroupId}

	// get all the threads
	if err := model.List(ctx, thread, &threads, model.Query{}); err != nil {
		logger.Errorf("Error reading thread: %v", err)
		return errors.InternalServerError("DATABASE_ERROR", "Error connecting to database")
	}

	// return the response
	for _, thread := range threads {
		rsp.Threads = append(rsp.Threads, thread.Serialize())
	}

	return nil
}
