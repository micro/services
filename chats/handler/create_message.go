package handler

import (
	"context"

	"github.com/google/uuid"
	"github.com/micro/micro/v3/service/errors"
	"github.com/micro/micro/v3/service/logger"
	pb "github.com/micro/services/chats/proto"
	"gorm.io/gorm"
)

// Create a message within a chat
func (c *Chats) CreateMessage(ctx context.Context, req *pb.CreateMessageRequest, rsp *pb.CreateMessageResponse) error {
	// validate the request
	if len(req.AuthorId) == 0 {
		return ErrMissingAuthorID
	}
	if len(req.ChatId) == 0 {
		return ErrMissingChatID
	}
	if len(req.Text) == 0 {
		return ErrMissingText
	}

	return c.DB.Transaction(func(tx *gorm.DB) error {
		// lookup the chat
		var conv Chat
		if err := tx.Where(&Chat{ID: req.ChatId}).First(&conv).Error; err == gorm.ErrRecordNotFound {
			return ErrNotFound
		} else if err != nil {
			logger.Errorf("Error reading chat: %v", err)
			return errors.InternalServerError("DATABASE_ERROR", "Error connecting to database")
		}

		// create the message
		msg := &Message{
			ID:       uuid.New().String(),
			SentAt:   c.Time(),
			Text:     req.Text,
			AuthorID: req.AuthorId,
			ChatID:   req.ChatId,
		}
		if err := tx.Create(msg).Error; err != nil {
			logger.Errorf("Error creating message: %v", err)
			return errors.InternalServerError("DATABASE_ERROR", "Error connecting to database")
		}

		// serialize the response
		rsp.Message = msg.Serialize()
		return nil
	})
}
