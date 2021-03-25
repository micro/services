package handler

import (
	"context"

	"github.com/micro/micro/v3/service/errors"
	"github.com/micro/micro/v3/service/logger"
	pb "github.com/micro/services/threads/proto"
	"gorm.io/gorm"
)

// RecentMessages returns the most recent messages in a group of conversations. By default the
// most messages retrieved per conversation is 25, however this can be overriden using the
// limit_per_conversation option
func (s *Threads) RecentMessages(ctx context.Context, req *pb.RecentMessagesRequest, rsp *pb.RecentMessagesResponse) error {
	// validate the request
	if len(req.ConversationIds) == 0 {
		return ErrMissingConversationIDs
	}

	limit := DefaultLimit
	if req.LimitPerConversation != nil {
		limit = int(req.LimitPerConversation.Value)
	}

	db, err := s.GetDBConn(ctx)
	if err != nil {
		logger.Errorf("Error connecting to DB: %v", err)
		return errors.InternalServerError("DB_ERROR", "Error connecting to DB")
	}
	// query the database
	var msgs []Message
	err = db.Transaction(func(tx *gorm.DB) error {
		for _, id := range req.ConversationIds {
			var cms []Message
			if err := tx.Where(&Message{ConversationID: id}).Order("sent_at DESC").Limit(limit).Find(&cms).Error; err != nil {
				return err
			}
			msgs = append(msgs, cms...)
		}
		return nil
	})
	if err != nil {
		logger.Errorf("Error reading messages: %v", err)
		return errors.InternalServerError("DATABASE_ERROR", "Error connecting to database")
	}

	// serialize the response
	rsp.Messages = make([]*pb.Message, len(msgs))
	for i, m := range msgs {
		rsp.Messages[i] = m.Serialize()
	}
	return nil
}
