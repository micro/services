package handler

import (
	"context"

	"github.com/micro/micro/v3/service/auth"
	"github.com/micro/micro/v3/service/errors"
	"github.com/micro/micro/v3/service/logger"
	"github.com/micro/services/pkg/model"
	pb "github.com/micro/services/threads/proto"
)

// Read a thread using its ID, can filter using group ID if provided
func (s *Threads) ReadThread(ctx context.Context, req *pb.ReadThreadRequest, rsp *pb.ReadThreadResponse) error {
	_, ok := auth.AccountFromContext(ctx)
	if !ok {
		errors.Unauthorized("UNAUTHORIZED", "Unauthorized")
	}
	// validate the request
	if len(req.Id) == 0 {
		return ErrMissingID
	}

	// construct the query
	thread := &Thread{ID: req.Id}

	var err error

	if len(req.GroupId) > 0 {
		thread.GroupID = req.GroupId
		err = model.ReadIndex(ctx, thread)
	} else {
		err = model.Read(ctx, thread)
	}

	if err == model.ErrNotFound {
		return ErrNotFound
	} else if err != nil {
		logger.Errorf("Error reading thread: %v", err)
		return errors.InternalServerError("DATABASE_ERROR", "Error connecting to database")
	}

	// serialize the response
	rsp.Thread = thread.Serialize()
	return nil
}
