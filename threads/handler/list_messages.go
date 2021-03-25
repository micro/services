package handler

import (
	"context"

	"github.com/micro/micro/v3/service/errors"
	"github.com/micro/micro/v3/service/logger"
	pb "github.com/micro/services/threads/proto"
)

const DefaultLimit = 25

// List the messages within a conversation in reverse chronological order, using sent_before to
// offset as older messages need to be loaded
func (s *Threads) ListMessages(ctx context.Context, req *pb.ListMessagesRequest, rsp *pb.ListMessagesResponse) error {
	// validate the request
	if len(req.ConversationId) == 0 {
		return ErrMissingConversationID
	}

	db, err := s.GetDBConn(ctx)
	if err != nil {
		logger.Errorf("Error connecting to DB: %v", err)
		return errors.InternalServerError("DB_ERROR", "Error connecting to DB")
	}
	// construct the query
	q := db.Where(&Message{ConversationID: req.ConversationId}).Order("sent_at DESC")
	if req.SentBefore != nil {
		q = q.Where("sent_at < ?", req.SentBefore.AsTime())
	}
	if req.Limit != nil {
		q.Limit(int(req.Limit.Value))
	} else {
		q.Limit(DefaultLimit)
	}

	// execute the query
	var msgs []Message
	if err := q.Find(&msgs).Error; err != nil {
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
