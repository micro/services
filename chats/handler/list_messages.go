package handler

import (
	"context"

	"github.com/micro/micro/v3/service/errors"
	"github.com/micro/micro/v3/service/logger"
	pb "github.com/micro/services/chats/proto"
)

const DefaultLimit = 25

// List the messages within a chat in reverse chronological order, using sent_before to
// offset as older messages need to be loaded
func (c *Chats) ListMessages(ctx context.Context, req *pb.ListMessagesRequest, rsp *pb.ListMessagesResponse) error {
	// validate the request
	if len(req.ChatId) == 0 {
		return ErrMissingChatID
	}

	// construct the query
	q := c.DB.Where(&Message{ChatID: req.ChatId}).Order("sent_at DESC")
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
