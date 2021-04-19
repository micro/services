package handler

import (
	"context"

	"github.com/micro/micro/v3/service/auth"
	"gorm.io/gorm"

	"github.com/micro/micro/v3/service/errors"
	"github.com/micro/micro/v3/service/logger"
	pb "github.com/micro/services/threads/proto"
)

// Read a conversation using its ID, can filter using group ID if provided
func (s *Threads) ReadConversation(ctx context.Context, req *pb.ReadConversationRequest, rsp *pb.ReadConversationResponse) error {
	_, ok := auth.AccountFromContext(ctx)
	if !ok {
		errors.Unauthorized("UNAUTHORIZED", "Unauthorized")
	}
	// validate the request
	if len(req.Id) == 0 {
		return ErrMissingID
	}

	// construct the query
	q := Conversation{ID: req.Id}
	if req.GroupId != "" {
		q.GroupID = req.GroupId
	}

	db, err := s.GetDBConn(ctx)
	if err != nil {
		logger.Errorf("Error connecting to DB: %v", err)
		return errors.InternalServerError("DB_ERROR", "Error connecting to DB")
	}
	// execute the query
	var conv Conversation
	if err := db.Where(&q).First(&conv).Error; err == gorm.ErrRecordNotFound {
		return ErrNotFound
	} else if err != nil {
		logger.Errorf("Error reading conversation: %v", err)
		return errors.InternalServerError("DATABASE_ERROR", "Error connecting to database")
	}

	// serialize the response
	rsp.Conversation = conv.Serialize()
	return nil
}
