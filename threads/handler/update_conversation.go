package handler

import (
	"context"

	"github.com/micro/micro/v3/service/errors"
	"github.com/micro/micro/v3/service/logger"
	pb "github.com/micro/services/threads/proto"
	"gorm.io/gorm"
)

// Update a conversations topic
func (s *Threads) UpdateConversation(ctx context.Context, req *pb.UpdateConversationRequest, rsp *pb.UpdateConversationResponse) error {
	// validate the request
	if len(req.Id) == 0 {
		return ErrMissingID
	}
	if len(req.Topic) == 0 {
		return ErrMissingTopic
	}

	// lookup the conversation
	var conv Conversation
	if err := s.DB.Where(&Conversation{ID: req.Id}).First(&conv).Error; err == gorm.ErrRecordNotFound {
		return ErrNotFound
	} else if err != nil {
		logger.Errorf("Error reading conversation: %v", err)
		return errors.InternalServerError("DATABASE_ERROR", "Error connecting to database")
	}

	// update the conversation
	conv.Topic = req.Topic
	if err := s.DB.Save(&conv).Error; err != nil {
		logger.Errorf("Error updating conversation: %v", err)
		return errors.InternalServerError("DATABASE_ERROR", "Error connecting to database")
	}

	// serialize the result
	rsp.Conversation = conv.Serialize()
	return nil
}
