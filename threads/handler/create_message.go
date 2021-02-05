package handler

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"github.com/micro/micro/v3/service/errors"
	"github.com/micro/micro/v3/service/logger"
	pb "github.com/micro/services/threads/proto"
	"gorm.io/gorm"
)

// Create a message within a conversation
func (s *Threads) CreateMessage(ctx context.Context, req *pb.CreateMessageRequest, rsp *pb.CreateMessageResponse) error {
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
		IdempotentID:   req.IdempotentId,
	}
	if len(msg.IdempotentID) == 0 {
		msg.IdempotentID = uuid.New().String()
	}
	if err := s.DB.Create(msg).Error; err == nil {
		rsp.Message = msg.Serialize()
		return nil
	} else if !strings.Contains(err.Error(), "idempotent_id") {
		logger.Errorf("Error creating message: %v", err)
		return errors.InternalServerError("DATABASE_ERROR", "Error connecting to database")
	}

	// a message already exists with this idempotent_id
	var existing Message
	if err := s.DB.Where(&Message{IdempotentID: msg.IdempotentID}).First(&existing).Error; err != nil {
		logger.Errorf("Error creating message: %v", err)
		return errors.InternalServerError("DATABASE_ERROR", "Error connecting to database")
	}
	rsp.Message = existing.Serialize()
	return nil
}
