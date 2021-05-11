package handler

import (
	"context"

	"github.com/micro/micro/v3/service/auth"
	"github.com/micro/micro/v3/service/errors"
	"github.com/micro/micro/v3/service/logger"
	"github.com/micro/micro/v3/service/store"
	pb "github.com/micro/services/threads/proto"
)

const DefaultLimit = 25

// List the messages within a thread in reverse chronological order, using sent_before to
// offset as older messages need to be loaded
func (s *Threads) ListMessages(ctx context.Context, req *pb.ListMessagesRequest, rsp *pb.ListMessagesResponse) error {
	_, ok := auth.AccountFromContext(ctx)
	if !ok {
		errors.Unauthorized("UNAUTHORIZED", "Unauthorized")
	}
	// validate the request
	if len(req.ThreadId) == 0 {
		return ErrMissingThreadID
	}

	// default order is descending
	order := store.OrderDesc
	if req.Order == "asc" {
		order = store.OrderAsc
	}

	opts := []store.ReadOption{
		store.ReadPrefix(),
		store.ReadOrder(order),
	}

	if req.Limit > 0 {
		opts = append(opts, store.ReadLimit(uint(req.Limit)))
	} else {
		opts = append(opts, store.ReadLimit(uint(DefaultLimit)))
	}
	if req.Offset > 0 {
		opts = append(opts, store.ReadOffset(uint(req.Offset)))
	}

	message := &Message{
		ThreadID: req.ThreadId,
	}

	// read all the records with the chat ID suffix
	recs, err := store.Read(message.Index(ctx), opts...)
	if err != nil {
		logger.Errorf("Error reading messages: %v", err)
		return errors.InternalServerError("DATABASE_ERROR", "Error connecting to database")
	}

	// return all the messages
	for _, rec := range recs {
		m := &Message{}
		rec.Decode(&m)
		if len(m.ID) == 0 || m.ThreadID != req.ThreadId {
			continue
		}
		rsp.Messages = append(rsp.Messages, m.Serialize())
	}

	return nil
}
