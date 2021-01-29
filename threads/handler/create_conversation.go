package handler

import (
	"context"

	"github.com/google/uuid"
	"github.com/micro/micro/v3/service/errors"
	"github.com/micro/micro/v3/service/logger"
	pb "github.com/micro/services/threads/proto"
)

// Create a conversation
func (s *Threads) CreateConversation(ctx context.Context, req *pb.CreateConversationRequest, rsp *pb.CreateConversationResponse) error {
	// validate the request
	if len(req.GroupId) == 0 {
		return ErrMissingGroupID
	}
	if len(req.Topic) == 0 {
		return ErrMissingTopic
	}

	// write the conversation to the database
	conv := &Conversation{
		ID:        uuid.New().String(),
		Topic:     req.Topic,
		GroupID:   req.GroupId,
		CreatedAt: s.Time(),
	}
	if err := s.DB.Create(conv).Error; err != nil {
		logger.Errorf("Error creating conversation: %v", err)
		return errors.InternalServerError("DATABASE_ERROR", "Error connecting to database")
	}

	// serialize the response
	rsp.Conversation = conv.Serialize()
	return nil
}
