package handler

import (
	"context"

	"github.com/micro/micro/v3/service/auth"
	"github.com/micro/micro/v3/service/errors"
	"github.com/micro/micro/v3/service/logger"
	"github.com/micro/services/pkg/model"
	pb "github.com/micro/services/threads/proto"
)

// Delete a thread and all the messages within
func (s *Threads) DeleteThread(ctx context.Context, req *pb.DeleteThreadRequest, rsp *pb.DeleteThreadResponse) error {
	_, ok := auth.AccountFromContext(ctx)
	if !ok {
		errors.Unauthorized("UNAUTHORIZED", "Unauthorized")
	}
	// validate the request
	if len(req.Id) == 0 {
		return ErrMissingID
	}

	thread := Thread{ID: req.Id}

	// delete the thread
	if err := model.Delete(ctx, &thread); err != nil {
		logger.Errorf("Error deleting thread: %v", err)
		return errors.InternalServerError("DATABASE_ERROR", "Error connecting to database")
	}

	message := Message{ThreadID: req.Id}

	// delete the messages
	if err := model.Delete(ctx, &message); err != nil {
		logger.Errorf("Error deleting messages: %v", err)
		return errors.InternalServerError("DATABASE_ERROR", "Error connecting to database")
	}

	return nil
}
