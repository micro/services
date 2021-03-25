package handler

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"github.com/micro/micro/v3/service/auth"
	"github.com/micro/micro/v3/service/errors"
	"github.com/micro/micro/v3/service/logger"
	pb "github.com/micro/services/threads/proto"
	"gorm.io/gorm"
)

// Create a message within a conversation
func (s *Threads) CreateMessage(ctx context.Context, req *pb.CreateMessageRequest, rsp *pb.CreateMessageResponse) error {
	_, ok := auth.AccountFromContext(ctx)
	if !ok {
		errors.Unauthorized("UNAUTHORIZED", "Unauthorized")
	}
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

	db, err := s.GetDBConn(ctx)
	if err != nil {
		logger.Errorf("Error connecting to DB: %v", err)
		return errors.InternalServerError("DB_ERROR", "Error connecting to DB")
	}
	// lookup the conversation
	var conv Conversation
	if err := db.Where(&Conversation{ID: req.ConversationId}).First(&conv).Error; err == gorm.ErrRecordNotFound {
		return ErrNotFound
	} else if err != nil {
		logger.Errorf("Error reading conversation: %v", err)
		return errors.InternalServerError("DATABASE_ERROR", "Error connecting to database")
	}

	// create the message
	msg := &Message{
		ID:             req.Id,
		SentAt:         s.Time(),
		Text:           req.Text,
		AuthorID:       req.AuthorId,
		ConversationID: req.ConversationId,
	}
	if len(msg.ID) == 0 {
		msg.ID = uuid.New().String()
	}
	if err := db.Create(msg).Error; err == nil {
		rsp.Message = msg.Serialize()
		return nil
	} else if !strings.Contains(err.Error(), "messages_pkey") {
		logger.Errorf("Error creating message: %v", err)
		return errors.InternalServerError("DATABASE_ERROR", "Error connecting to database")
	}

	// a message already exists with this id
	var existing Message
	if err := db.Where(&Message{ID: msg.ID}).First(&existing).Error; err != nil {
		logger.Errorf("Error creating message: %v", err)
		return errors.InternalServerError("DATABASE_ERROR", "Error connecting to database")
	}
	rsp.Message = existing.Serialize()
	return nil
}
