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

// Create a thread
func (s *Threads) CreateThread(ctx context.Context, req *pb.CreateThreadRequest, rsp *pb.CreateThreadResponse) error {
	_, ok := auth.AccountFromContext(ctx)
	if !ok {
		errors.Unauthorized("UNAUTHORIZED", "Unauthorized")
	}
	// validate the request
	if len(req.GroupId) == 0 {
		return ErrMissingGroupID
	}
	if len(req.Topic) == 0 {
		return ErrMissingTopic
	}

	// write the thread to the database
	thread := &Thread{
		ID:        uuid.New().String(),
		Topic:     req.Topic,
		GroupID:   req.GroupId,
		CreatedAt: s.Time(),
	}

	// write the thread to the database
	if err := model.Create(ctx, thread); err != nil {
		logger.Errorf("Error creating thread: %v", err)
		return err
	}

	// serialize the response
	rsp.Thread = thread.Serialize()
	return nil
}
