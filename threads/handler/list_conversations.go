package handler

import (
	"context"

	"github.com/micro/micro/v3/service/errors"
	"github.com/micro/micro/v3/service/logger"
	pb "github.com/micro/services/threads/proto"
)

// List all the conversations for a group
func (s *Threads) ListConversations(ctx context.Context, req *pb.ListConversationsRequest, rsp *pb.ListConversationsResponse) error {
	// validate the request
	if len(req.GroupId) == 0 {
		return ErrMissingGroupID
	}

	db, err := s.GetDBConn(ctx)
	if err != nil {
		logger.Errorf("Error connecting to DB: %v", err)
		return errors.InternalServerError("DB_ERROR", "Error connecting to DB")
	}
	// query the database
	var convs []Conversation
	if err := db.Where(&Conversation{GroupID: req.GroupId}).Find(&convs).Error; err != nil {
		logger.Errorf("Error reading conversation: %v", err)
		return errors.InternalServerError("DATABASE_ERROR", "Error connecting to database")
	}

	// serialize the response
	rsp.Conversations = make([]*pb.Conversation, len(convs))
	for i, c := range convs {
		rsp.Conversations[i] = c.Serialize()
	}
	return nil
}
