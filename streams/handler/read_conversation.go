package handler

import (
	"context"

	"gorm.io/gorm"

	"github.com/micro/micro/v3/service/errors"
	"github.com/micro/micro/v3/service/logger"
	pb "github.com/micro/services/streams/proto"
)

// Read a conversation using its ID, can filter using group ID if provided
func (s *Streams) ReadConversation(ctx context.Context, req *pb.ReadConversationRequest, rsp *pb.ReadConversationResponse) error {
	// validate the request
	if len(req.Id) == 0 {
		return ErrMissingID
	}

	// construct the query
	q := Conversation{ID: req.Id}
	if req.GroupId != nil {
		q.GroupID = req.GroupId.Value
	}

	// execute the query
	var conv Conversation
	if err := s.DB.Where(&q).First(&conv).Error; err == gorm.ErrRecordNotFound {
		return ErrNotFound
	} else if err != nil {
		logger.Errorf("Error reading conversation: %v", err)
		return errors.InternalServerError("DATABASE_ERROR", "Error connecting to database")
	}

	// serialize the response
	rsp.Conversation = conv.Serialize()
	return nil
}
