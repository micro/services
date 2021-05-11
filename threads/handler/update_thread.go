package handler

import (
	"context"

	"github.com/micro/micro/v3/service/auth"
	"github.com/micro/micro/v3/service/errors"
	"github.com/micro/micro/v3/service/logger"
	"github.com/micro/services/pkg/model"
	pb "github.com/micro/services/threads/proto"
)

// Update a threads topic
func (s *Threads) UpdateThread(ctx context.Context, req *pb.UpdateThreadRequest, rsp *pb.UpdateThreadResponse) error {
	_, ok := auth.AccountFromContext(ctx)
	if !ok {
		errors.Unauthorized("UNAUTHORIZED", "Unauthorized")
	}
	// validate the request
	if len(req.Id) == 0 {
		return ErrMissingID
	}
	if len(req.Topic) == 0 {
		return ErrMissingTopic
	}

	t := &Thread{ID: req.Id}

	if err := model.Read(ctx, t); err == model.ErrNotFound {
		return ErrNotFound
	} else if err != nil {
		logger.Errorf("Error reading thread: %v", err)
		return errors.InternalServerError("DATABASE_ERROR", "Error connecting to database")
	}

	// update the thread
	t.Topic = req.Topic
	if err := model.Update(ctx, t); err != nil {
		logger.Errorf("Error updating thread: %v", err)
		return errors.InternalServerError("DATABASE_ERROR", "Error connecting to database")
	}

	// serialize the result
	rsp.Thread = t.Serialize()

	return nil
}
