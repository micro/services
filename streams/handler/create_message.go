package handler

import (
	"context"

	"github.com/google/uuid"
	"github.com/micro/micro/v3/service/errors"
	"github.com/micro/micro/v3/service/logger"
	pb "github.com/micro/services/streams/proto"
	"gorm.io/gorm"
)

// Create a message within a conversation
func (s *Streams) CreateMessage(ctx context.Context, req *pb.CreateMessageRequest, rsp *pb.CreateMessageResponse) error {
	// validate the request
	if len(req.AuthorId) == 0 {
		return ErrMissingAuthorID
	}
	if len(req.ConversationId) == 0 {
		return ErrMissingConversationID
	}
	if len(req.Text) == 0 {
		return ErrMissingText
	}

	return s.DB.Transaction(func(tx *gorm.DB) error {
		// lookup the conversation
		var conv Conversation
		if err := s.DB.Where(&Conversation{ID: req.ConversationId}).First(&conv).Error; err == gorm.ErrRecordNotFound {
			return ErrNotFound
		} else if err != nil {
			logger.Errorf("Error reading conversation: %v", err)
			return errors.InternalServerError("DATABASE_ERROR", "Error connecting to database")
		}

		// create the message
		msg := &Message{
			ID:             uuid.New().String(),
			SentAt:         s.Time(),
			Text:           req.Text,
			AuthorID:       req.AuthorId,
			ConversationID: req.ConversationId,
		}
		if err := s.DB.Create(msg).Error; err != nil {
			logger.Errorf("Error creating message: %v", err)
			return errors.InternalServerError("DATABASE_ERROR", "Error connecting to database")
		}

		// serialize the response
		rsp.Message = msg.Serialize()
		return nil
	})
}
