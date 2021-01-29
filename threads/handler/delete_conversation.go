package handler

import (
	"context"

	"github.com/micro/micro/v3/service/errors"
	"github.com/micro/micro/v3/service/logger"
	pb "github.com/micro/services/threads/proto"
	"gorm.io/gorm"
)

// Delete a conversation and all the messages within
func (s *Threads) DeleteConversation(ctx context.Context, req *pb.DeleteConversationRequest, rsp *pb.DeleteConversationResponse) error {
	// validate the request
	if len(req.Id) == 0 {
		return ErrMissingID
	}

	return s.DB.Transaction(func(tx *gorm.DB) error {
		// delete all the messages
		if err := tx.Where(&Message{ConversationID: req.Id}).Delete(&Message{}).Error; err != nil {
			logger.Errorf("Error deleting messages: %v", err)
			return errors.InternalServerError("DATABASE_ERROR", "Error connecting to database")
		}

		// delete the conversation
		if err := tx.Where(&Conversation{ID: req.Id}).Delete(&Conversation{}).Error; err != nil {
			logger.Errorf("Error deleting conversation: %v", err)
			return errors.InternalServerError("DATABASE_ERROR", "Error connecting to database")
		}

		return nil
	})
}
