package handler

import (
	"context"

	"github.com/micro/micro/v3/service/auth"
	"github.com/micro/micro/v3/service/errors"
	"github.com/micro/micro/v3/service/logger"
	"github.com/micro/services/pkg/model"
	pb "github.com/micro/services/threads/proto"
)

// RecentMessages returns the most recent messages in a group of threads. By default the
// most messages retrieved per thread is 25, however this can be overriden using the
// limit_per_thread option
func (s *Threads) RecentMessages(ctx context.Context, req *pb.RecentMessagesRequest, rsp *pb.RecentMessagesResponse) error {
	_, ok := auth.AccountFromContext(ctx)
	if !ok {
		errors.Unauthorized("UNAUTHORIZED", "Unauthorized")
	}
	// validate the request
	if len(req.ThreadIds) == 0 {
		return ErrMissingThreadIDs
	}

	limit := uint(DefaultLimit)
	if req.LimitPerThread > 0 {
		limit = uint(req.LimitPerThread)
	}

	for _, thread := range req.ThreadIds {
		q := model.Query{Limit: limit, Order: "desc"}
		m := &Message{ThreadID: thread}
		var messages []*Message

		if err := model.List(ctx, m, &messages, q); err != nil {
			logger.Errorf("Error reading messages: %v", err)
			return errors.InternalServerError("DATABASE_ERROR", "Error connecting to database")
		}

		for _, msg := range messages {
			rsp.Messages = append(rsp.Messages, msg.Serialize())
		}
	}

	return nil
}
