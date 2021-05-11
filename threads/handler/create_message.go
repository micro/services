package handler

import (
	"context"

	"github.com/google/uuid"
	"github.com/micro/micro/v3/service/auth"
	"github.com/micro/micro/v3/service/errors"
	"github.com/micro/micro/v3/service/logger"
	"github.com/micro/services/pkg/model"
	pb "github.com/micro/services/threads/proto"
)

// Create a message within a thread
func (s *Threads) CreateMessage(ctx context.Context, req *pb.CreateMessageRequest, rsp *pb.CreateMessageResponse) error {
	_, ok := auth.AccountFromContext(ctx)
	if !ok {
		errors.Unauthorized("UNAUTHORIZED", "Unauthorized")
	}
	// validate the request
	if len(req.AuthorId) == 0 {
		return ErrMissingAuthorID
	}
	if len(req.ThreadId) == 0 {
		return ErrMissingThreadID
	}
	if len(req.Text) == 0 {
		return ErrMissingText
	}

	// lookup the thread
	conv := Thread{ID: req.ThreadId}

	if err := model.Read(ctx, &conv); err == model.ErrNotFound {
		return ErrNotFound
	} else if err != nil {
		logger.Errorf("Error reading thread: %v", err)
		return errors.InternalServerError("DATABASE_ERROR", "Error connecting to database")
	}

	// create the message
	msg := &Message{
		ID:       req.Id,
		SentAt:   s.Time(),
		Text:     req.Text,
		AuthorID: req.AuthorId,
		ThreadID: req.ThreadId,
	}
	if len(msg.ID) == 0 {
		msg.ID = uuid.New().String()
	}

	if err := model.Create(ctx, msg); err == nil {
		rsp.Message = msg.Serialize()
		return nil
	} else if err != model.ErrAlreadyExists {
		logger.Errorf("Error creating message: %v", err)
		return errors.InternalServerError("DATABASE_ERROR", "Error connecting to database")
	}

	// a message already exists with this id
	existing := &Message{ID: msg.ID, ThreadID: req.ThreadId}

	if err := model.Read(ctx, existing); err == model.ErrNotFound {
		return ErrNotFound
	} else if err != nil {
		logger.Errorf("Error creating message: %v", err)
		return errors.InternalServerError("DATABASE_ERROR", "Error connecting to database")
	}

	// return the message
	rsp.Message = existing.Serialize()
	return nil
}
